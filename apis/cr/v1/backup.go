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

const PgBackupResourcePlural = "pgbackups"

type PgBackupSpec struct {
	Name          string        `json:"name"`
	StorageSpec   PgStorageSpec `json:"storagespec"`
	CCP_IMAGE_TAG string        `json:"ccpimagetag"`
	BACKUP_HOST   string        `json:"backuphost"`
	BACKUP_USER   string        `json:"backupuser"`
	BACKUP_PASS   string        `json:"backuppass"`
	BACKUP_PORT   string        `json:"backupport"`
	BACKUP_STATUS string        `json:"backupstatus"`
}

type PgBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   PgBackupSpec   `json:"spec"`
	Status PgBackupStatus `json:"status,omitempty"`
}

type PgBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []PgBackup `json:"items"`
}

type PgBackupStatus struct {
	State   PgBackupState `json:"state,omitempty"`
	Message string        `json:"message,omitempty"`
}

type PgBackupState string

const (
	PgBackupStateCreated   PgBackupState = "Created"
	PgBackupStateProcessed PgBackupState = "Processed"
)
