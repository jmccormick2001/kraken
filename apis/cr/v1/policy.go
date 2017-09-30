/*
 Copyright 2017 Crunchy Data Solutions, Inc.
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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const PgPolicyResourcePlural = "pgpolicies"

type PgPolicySpec struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Sql    string `json:"sql"`
	Status string `json:"status"`
}

type PgPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   PgPolicySpec   `json:"spec"`
	Status PgPolicyStatus `json:"status,omitempty"`
}

type PgPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PgPolicy `json:"items"`
}

type PgPolicyStatus struct {
	State   PgPolicyState `json:"state,omitempty"`
	Message string        `json:"message,omitempty"`
}

type PgPolicyState string

const (
	PgPolicyStateCreated   PgPolicyState = "Created"
	PgPolicyStateProcessed PgPolicyState = "Processed"
)
