package resource

import (
	"fmt"
	"testing"

	"github.com/weaveworks/flux/resource"
)

func TestParseImageOnlyFormat(t *testing.T) {
	expectedImage := "bitnami/mariadb:10.1.30-r1"
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    image: ` + expectedImage + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
}

func TestParseImageTagFormat(t *testing.T) {
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    image: ` + expectedImageName + `
    tag: ` + expectedImageTag + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
}

func TestParseNamedImageFormat(t *testing.T) {
	expectedContainer := "db"
	expectedImage := "bitnami/mariadb:10.1.30-r1"
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    ` + expectedContainer + `:
      first: post
      image: ` + expectedImage + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseNamedImageTagFormat(t *testing.T) {
	expectedContainer := "db"
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    other:
      not: "containing image"
    ` + expectedContainer + `:
      first: post
      image: ` + expectedImageName + `
      tag: ` + expectedImageTag + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseImageObjectFormat(t *testing.T) {
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    image:
      repository: ` + expectedImageName + `
      tag: ` + expectedImageTag + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
}

func TestParseNamedImageObjectFormat(t *testing.T) {
	expectedContainer := "db"
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    other:
      not: "containing image"
    ` + expectedContainer + `:
      first: post
      image:
        repository: ` + expectedImageName + `
        tag: ` + expectedImageTag + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseAllFormatsInOne(t *testing.T) {

	type container struct {
		name, image, tag string
	}

	// *NB* the containers will be calculated based on the order
	//  1. the entry for 'image' if present
	//  2. the order of the keys in `values`.
	//
	// To avoid having to mess around later, I have cooked the order
	// of these so they can be compared directly to the return value.
	expected := []container{
		{ReleaseContainerName, "repo/imageOne", "tagOne"},
		{"AAA", "repo/imageTwo", "tagTwo"},
		{"ZZZ", "repo/imageThree", "tagThree"},
	}

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: test
  namespace: test
spec:
  chartGitPath: test
  values:
    # top-level image
    image: ` + expected[0].image + ":" + expected[0].tag + `

    # under .container, as image and tag entries
    ` + expected[1].name + `:
      image: ` + expected[1].image + `
      tag: ` + expected[1].tag + `

    # under .container.image, as repository and tag entries
    ` + expected[2].name + `:
      image:
        repository: ` + expected[2].image + `
        tag: ` + expected[2].tag + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["test:fluxhelmrelease/test"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != len(expected) {
		t.Fatalf("expected %d containers, got %d", len(expected), len(containers))
	}
	for i, c0 := range expected {
		c1 := containers[i]
		if c1.Name != c0.name {
			t.Errorf("names do not match %q != %q", c0, c1)
		}
		c0image := fmt.Sprintf("%s:%s", c0.image, c0.tag)
		if c1.Image.String() != c0image {
			t.Errorf("images do not match %q != %q", c0image, c1.Image.String())
		}
	}
}
