package main

import (
	"context"
	"os"

	eraserv1alpha1 "github.com/Azure/eraser/api/v1alpha1"
	util "github.com/Azure/eraser/pkg/utils"
	"github.com/davecgh/go-spew/spew"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	"k8s.io/client-go/tools/reference"
	"k8s.io/kubectl/pkg/scheme"
)

func removeImages(c Client, targetImages []string, recorder events.EventRecorder) error {
	backgroundContext, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	images, err := c.listImages(backgroundContext)
	if err != nil {
		return err
	}

	allImages := make([]string, 0, len(images))

	// map with key: sha id, value: repoTag list (contains full name of image)
	idToTagListMap := make(map[string][]string)

	for _, img := range images {
		allImages = append(allImages, img.Id)
		idToTagListMap[img.Id] = img.RepoTags
	}

	containers, err := c.listContainers(backgroundContext)
	if err != nil {
		return err
	}

	// Images that are running
	// map of (digest | tag) -> digest
	runningImages := util.GetRunningImages(containers, idToTagListMap)

	// Images that aren't running
	// map of (digest | tag) -> digest
	nonRunningImages := util.GetNonRunningImages(runningImages, allImages, idToTagListMap)

	// Debug logs
	log.V(1).Info("Map of non-running images", "nonRunningImages", nonRunningImages)
	log.V(1).Info("Map of running images", "runningImages", runningImages)
	log.V(1).Info("Map of digest to image name(s)", "idToTaglistMap", idToTagListMap)

	// remove target images
	var prune bool
	deletedImages := make(map[string]struct{}, len(targetImages))
	for _, imgDigestOrTag := range targetImages {
		if imgDigestOrTag == "*" {
			prune = true
			continue
		}

		if digest, isNonRunning := nonRunningImages[imgDigestOrTag]; isNonRunning {
			if ex := util.IsExcluded(excluded, imgDigestOrTag, idToTagListMap); ex {
				log.Info("image is excluded", "given", imgDigestOrTag, "digest", digest, "name", idToTagListMap[digest])
				continue
			}

			// in the case name and digest are both provided
			if _, deleted := deletedImages[digest]; deleted {
				log.Info("image with digest already deleted", "given", imgDigestOrTag, "digest", digest, "name", idToTagListMap[digest])
				continue
			}

			err = c.deleteImage(backgroundContext, digest)
			if err != nil {
				log.Error(err, "error removing image", "given", imgDigestOrTag, "digest", digest, "name", idToTagListMap[digest])
				continue
			}

			log.Info("removed image", "given", imgDigestOrTag, "digest", digest, "name", idToTagListMap[digest])
			deletedImages[imgDigestOrTag] = struct{}{}
		}

		digest, isRunning := runningImages[imgDigestOrTag]
		if isRunning {
			log.Info("image is running", "given", imgDigestOrTag, "digest", digest, "name", idToTagListMap[digest])
			continue
		}
		log.Info("image is not on node", "given", imgDigestOrTag)
	}

	if prune {
		success := true
		for _, digest := range nonRunningImages {
			// in the case both name -> digest and digest -> digest were present
			if _, deleted := deletedImages[digest]; deleted {
				continue
			}

			if util.IsExcluded(excluded, digest, idToTagListMap) {
				log.Info("image is excluded", "digest", digest, "name", idToTagListMap[digest])
				continue
			}

			if err := c.deleteImage(backgroundContext, digest); err != nil {
				success = false
				log.Error(err, "error removing image", "digest", digest, "name", idToTagListMap[digest])
				continue
			}
			log.Info("removed image", "digest", digest, "name", idToTagListMap[digest])
			deletedImages[digest] = struct{}{}
		}
		if success {
			log.Info("prune successful")
		} else {
			log.Info("error during prune")
		}
	}

	if *emitEvents {
		var finalRemoved []eraserv1alpha1.Image
		for digest := range deletedImages {
			if len(idToTagListMap[digest]) > 0 {
				finalRemoved = append(finalRemoved, eraserv1alpha1.Image{Digest: digest, Name: idToTagListMap[digest][0]})
			} else {
				finalRemoved = append(finalRemoved, eraserv1alpha1.Image{Digest: digest})
			}
		}
		emitEvent(recorder, finalRemoved)
	}

	return nil
}

func emitEvent(recorder events.EventRecorder, finalRemoved []eraserv1alpha1.Image) {
	nodeName := os.Getenv("NODE_NAME")
	nodeRef := &corev1.ObjectReference{
		Kind:      "Node",
		Name:      nodeName,
		UID:       types.UID(nodeName),
		Namespace: "",
	}

	spew.Dump(nodeRef)

	ref, err := reference.GetReference(scheme.Scheme, nodeRef)
	if err != nil {
		log.Error(err, "could not get reference to node", nodeName)
	}
	log.Info("recording event 1", "finalRemoved", finalRemoved)
	recorder.Eventf(ref, nil, corev1.EventTypeNormal, "RemovedImage", "removed image", "testing with short note %s", "test")
}
