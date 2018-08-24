package kubeapi

/*
 Copyright 2017-2018 Crunchy Data Solutions, Inc.
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

import (
	log "github.com/Sirupsen/logrus"
	crv1 "github.com/crunchydata/postgres-operator/apis/cr/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

// GetpgtasksBySelector gets a list of pgtasks by selector
func GetpgtasksBySelector(client *rest.RESTClient, taskList *crv1.PgtaskList, selector, namespace string) error {

	var err error

	myselector := labels.Everything()

	if selector != "" {
		myselector, err = labels.Parse(selector)
		if err != nil {
			log.Error("could not parse selector value ")
			log.Error(err)
			return err
		}
	}

	log.Debug("myselector is " + myselector.String())

	err = client.Get().
		Resource(crv1.PgtaskResourcePlural).
		Namespace(namespace).
		Param("labelSelector", myselector.String()).
		Do().
		Into(taskList)
	if err != nil {
		log.Error("error getting list of tasks " + err.Error())
	}

	return err
}

// Getpgtasks gets a list of pgtasks
func Getpgtasks(client *rest.RESTClient, taskList *crv1.PgtaskList, namespace string) error {

	err := client.Get().
		Resource(crv1.PgtaskResourcePlural).
		Namespace(namespace).
		Do().Into(taskList)
	if err != nil {
		log.Error("error getting list of tasks " + err.Error())
		return err
	}

	return err
}

// Getpgtask gets a pgtask by name
func Getpgtask(client *rest.RESTClient, task *crv1.Pgtask, name, namespace string) (bool, error) {

	err := client.Get().
		Resource(crv1.PgtaskResourcePlural).
		Namespace(namespace).
		Name(name).
		Do().Into(task)
	if kerrors.IsNotFound(err) {
		log.Debug("task " + name + " not found")
		return false, err
	}
	if err != nil {
		log.Error("error getting task " + err.Error())
		return true, err
	}

	return true, err
}

// Deletepgtask deletes pgtask by name
func Deletepgtask(client *rest.RESTClient, name, namespace string) error {

	err := client.Delete().
		Resource(crv1.PgtaskResourcePlural).
		Namespace(namespace).
		Name(name).
		Do().
		Error()
	if err != nil {
		log.Error("error deleting pgtask " + err.Error())
		return err
	}

	return err
}

// Createpgtask creates a pgtask
func Createpgtask(client *rest.RESTClient, task *crv1.Pgtask, namespace string) error {

	result := crv1.Pgtask{}

	err := client.Post().
		Resource(crv1.PgtaskResourcePlural).
		Namespace(namespace).
		Body(task).
		Do().
		Into(&result)
	if err != nil {
		log.Error("error creating pgtask " + err.Error())
	}

	log.Debug("created pgtask " + task.Name)
	return err
}

// Updatepgtask updates a pgtask
func Updatepgtask(client *rest.RESTClient, task *crv1.Pgtask, name, namespace string) error {

	err := client.Put().
		Name(name).
		Namespace(namespace).
		Resource(crv1.PgtaskResourcePlural).
		Body(task).
		Do().
		Error()
	if err != nil {
		log.Error("error updating pgtask " + err.Error())
	}

	log.Debug("updated pgtask " + task.Name)
	return err
}

// Deletepgtasks deletes pgtask by selector
func Deletepgtasks(client *rest.RESTClient, selector, namespace string) error {
	taskList := crv1.PgtaskList{}
	err := GetpgtasksBySelector(client, &taskList, selector, namespace)
	if err != nil {
		log.Error(err)
		return err
	}
	for _, v := range taskList.Items {
		err := Deletepgtask(client, v.ObjectMeta.Name, namespace)
		if err != nil {
			return err
		}
	}
	return err
}
