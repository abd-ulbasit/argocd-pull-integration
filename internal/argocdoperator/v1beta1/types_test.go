/*
Copyright 2026 Open Cluster Management.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

// operatorCRDFiles are the argocd-operator CRD bundles this repo installs.
var operatorCRDFiles = []struct {
	name string
	path string
}{
	{name: "hub chart", path: filepath.Join("..", "..", "..", "charts", "argocd-agent-addon", "crds", "argocd-operator-crds.yaml")},
	{name: "addon chart", path: filepath.Join("..", "..", "addon", "charts", "argocd-agent-addon", "crds", "argocd-operator-crds.yaml")},
}

// TestTypesMatchOperatorCRD checks every field of the ArgoCD type against the ArgoCD CRD of the
// operator version this repo ships. It fails when a field is missing from the CRD or has a
// different type, for example after an operator bump renames or removes a field.
func TestTypesMatchOperatorCRD(t *testing.T) {
	for _, file := range operatorCRDFiles {
		t.Run(file.name, func(t *testing.T) {
			schema := loadArgoCDSchema(t, file.path)
			checkType(t, ArgoCDKind, reflect.TypeOf(ArgoCD{}), schema)
		})
	}
}

// loadArgoCDSchema returns the OpenAPI schema of the ArgoCD CRD version this package mirrors.
func loadArgoCDSchema(t *testing.T, file string) *apiextensionsv1.JSONSchemaProps {
	t.Helper()

	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read %s: %v", file, err)
	}

	for _, doc := range strings.Split(string(data), "\n---\n") {
		if !strings.Contains(doc, "name: argocds."+GroupVersion.Group) {
			continue
		}
		crd := &apiextensionsv1.CustomResourceDefinition{}
		if err := yaml.Unmarshal([]byte(doc), crd); err != nil {
			t.Fatalf("failed to parse the ArgoCD CRD in %s: %v", file, err)
		}
		if crd.Spec.Group != GroupVersion.Group || crd.Spec.Names.Kind != ArgoCDKind {
			continue
		}
		for _, version := range crd.Spec.Versions {
			if version.Name != GroupVersion.Version {
				continue
			}
			if !version.Served {
				t.Fatalf("%s: version %s of the ArgoCD CRD is not served", file, GroupVersion.Version)
			}
			if version.Schema == nil || version.Schema.OpenAPIV3Schema == nil {
				t.Fatalf("%s: version %s of the ArgoCD CRD has no schema", file, GroupVersion.Version)
			}
			return version.Schema.OpenAPIV3Schema
		}
	}

	t.Fatalf("%s: no %s %s CRD found", file, GroupVersion, ArgoCDKind)
	return nil
}

// checkType checks that typ, as encoded by its JSON tags, fits schema.
func checkType(t *testing.T, path string, typ reflect.Type, schema *apiextensionsv1.JSONSchemaProps) {
	t.Helper()

	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		expectSchemaType(t, path, schema, "string")
	case reflect.Bool:
		expectSchemaType(t, path, schema, "boolean")
	case reflect.Slice:
		if !expectSchemaType(t, path, schema, "array") {
			return
		}
		if schema.Items == nil || schema.Items.Schema == nil {
			t.Errorf("%s: the CRD declares no item schema", path)
			return
		}
		checkType(t, path+"[]", typ.Elem(), schema.Items.Schema)
	case reflect.Struct:
		if !expectSchemaType(t, path, schema, "object") {
			return
		}
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
			// TypeMeta and ObjectMeta are the standard Kubernetes fields, not part of the CRD schema.
			if field.Anonymous {
				continue
			}
			if name == "" || name == "-" {
				t.Errorf("%s: field %s has no JSON name", path, field.Name)
				continue
			}
			property, ok := schema.Properties[name]
			if !ok {
				t.Errorf("%s.%s: not in the operator CRD", path, name)
				continue
			}
			checkType(t, path+"."+name, field.Type, &property)
		}
	default:
		t.Errorf("%s: unsupported Go kind %s", path, typ.Kind())
	}
}

func expectSchemaType(t *testing.T, path string, schema *apiextensionsv1.JSONSchemaProps, want string) bool {
	t.Helper()
	if schema.Type != want {
		t.Errorf("%s: the CRD type is %q, the Go type encodes %q", path, schema.Type, want)
		return false
	}
	return true
}
