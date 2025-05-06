# Kubernetes Advanced Assignment

**`Note:`** I have installed Kubernetes and minikube on my windows machine. (Refer Installation_kubernetes.docx`) 

Start Minikube using the following command:
<pre>minikube start --driver=docker </pre>

## Overview 
<pre> 📁 Kubernetes_Assignment_Advanced/
       ├── 📁 Helm Charts/
            ├── 📁 userdataapp_backend_chart/
            ├── 📁 userdataapp_db_chart/
       ├── 📁 Postman_file/
       ├── 📁 UserDataApp_Codebase/
            ├── 📁 backend/
            ├── 📁 database/
</pre>

📂 Directory Descriptions:
- **`Helm Charts/`**

    Contains Helm charts for deploying Kubernetes resources:

    - **`userdataapp_backend_chart/`** – Helm chart for the backend (server) application deployment.

    - **`userdataapp_db_chart/`** – Helm chart for deploying the database service and resources. (Backend application connects to this db service.)

- **`Postman_file/`**

    Contains Postman collection(s) for testing the application's APIs.

    You just need to import it in Postman and update the service port (assigned by minikube, when you run the service using minikube service command, explained later) to test the application's API.

- **`UserDataApp_Codebase/`**

    Contains full source code and database scripts for docker build image.

    - **`backend/`** – contains business logic, API handlers, Dockerfile and configs.

    - **`database/`** – contains scripts for DockerFile.

## Setup and run the project 

- Run `Windows PowerShell` and go to `Helm Charts/` directory (replace path below according to your placement of the project locally):
    <pre> >cd "C:\Kubernetes_Assignment_Advanced\Helm Charts" </pre>

- Create namespace:
    <pre> >kubectl create ns userdataapp-dev 
    namespace/userdataapp-dev created</pre>

### Deploy Application using helm charts:

- Deploy DB service and resources in Kubernetes cluster using helm chart:
    <pre> >helm install userdataapp-db ".\userdataapp_db_chart" -n userdataapp-dev

    NAME: userdataapp-db
    LAST DEPLOYED: Mon May  5 10:15:39 2025
    NAMESPACE: userdataapp-dev
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None</pre>

- Check db resources:
    <pre> >kubectl get all -n userdataapp-dev 

    NAME                                  READY   STATUS    RESTARTS   AGE
    pod/userdataapp-db-75955c9759-q95cl   1/1     Running   0          3m6s

    NAME                             TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
    service/service-userdataapp-db   ClusterIP   None         <<none>none>        1433/TCP   3m6s

    NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/userdataapp-db   1/1     1            1           3m6s

    NAME                                        DESIRED   CURRENT   READY   AGE
    replicaset.apps/userdataapp-db-75955c9759   1         1         1       3m6s
    </pre>

- Deploy backend application's resources in Kubernetes cluster using helm chart:
    <pre> >helm install userdataapp-backend .\userdataapp_backend_chart\ -n userdataapp-dev --set environment.env_dev=true

    NAME: userdataapp-backend
    LAST DEPLOYED: Mon May  5 10:24:40 2025
    NAMESPACE: userdataapp-dev
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None </pre>

- Check the backend application resources:

    <pre> >kubectl get all -n userdataapp-dev -l app=userdataapp-backend

    NAME                                       READY   STATUS    RESTARTS   AGE
    pod/userdataapp-backend-5f6d64dfc5-kcj4g   1/1     Running   0          21m

    NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/userdataapp-backend   1/1     1            1           21m

    NAME                                             DESIRED   CURRENT   READY   AGE
    replicaset.apps/userdataapp-backend-5f6d64dfc5   1         1         1       21m
    </pre>

    Sometimes we see that backend service is not in Ready state (it will show 0 under Ready, if not skip and move to next steps). 

    This could be because backend application's readiness probe is failing. Can be checked using below command (pick the pod name from previous command and change it in below command):
    <pre> >kubecrl describe po/userdataapp-backend-5f6d64dfc5-kcj4g -n userdataapp-dev
    Name:             userdataapp-backend-5f6d64dfc5-kcj4g
    Namespace:        userdataapp-dev
    Priority:         0
    ------
    ------
    ------
    
    Events:
    Type     Reason     Age                 From               Message
    ----     ------     ----                ----               -------
    Normal   Scheduled  27m                 default-scheduler  Successfully assigned userdataapp-dev/userdataapp-backend-5f6d64dfc5-kcj4g to minikube
    Normal   Pulled     27m                 kubelet            Container image "rajatjain20/userdataapp_backend:1.4" already present on machine
    Normal   Created    27m                 kubelet            Created container: userdataapp-backend
    Normal   Started    27m                 kubelet            Started container userdataapp-backend
    Warning  Unhealthy  12m (x52 over 27m)  kubelet            Readiness probe failed: HTTP probe failed with statuscode: 500 </pre>   

    The last statement shows that readiness probe is failed.
    This is because Readiness endpoint responds status code 500, if it is unable to connect to the db service.

    To make it work. we need to execute script manually:

    Get the pod name: 
    <pre> >kubectl get pods -n userdataapp-dev -l app=userdataapp-db 

    NAME                              READY   STATUS    RESTARTS   AGE
    userdataapp-db-75955c9759-q95cl   1/1     Running   0          58m</pre>

    Open an interactive Bash shell inside a running Kubernetes pod and then execute the script setup.sql:
    <pre> >kubectl exec -it po/userdataapp-db-75955c9759-q95cl -n userdataapp-dev -- /bin/bash </pre>
    <pre> # /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P admin@123 -d master -i setup.sql -C
    Changed database context to 'USERDATA'. </pre>
    <pre> # exit </pre>

    Now if we check the backend resources again, it will show Ready state to be 1.

- Check the services created by our helm charts and run backend service of type Nodeport:
    <pre> >kubectl get svc -n userdataapp-dev 

    NAME                             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
    service-userdataapp-backend      ClusterIP   10.110.18.46     <none>        3030/TCP         61m
    service-userdataapp-backend-np   NodePort    10.103.194.209   <none>        3031:30198/TCP   61m
    service-userdataapp-db           ClusterIP   None             <none>        1433/TCP         70m
    </pre> 
    <pre> >minikube service service-userdataapp-backend-np -n userdataapp-dev

    |-----------------|--------------------------------|-------------|---------------------------|
    |    NAMESPACE    |              NAME              | TARGET PORT |            URL            |
    |-----------------|--------------------------------|-------------|---------------------------|
    | userdataapp-dev | service-userdataapp-backend-np |        3031 | http://192.168.49.2:30198 |
    |-----------------|--------------------------------|-------------|---------------------------|
    * Starting tunnel for service service-userdataapp-backend-np.
    |-----------------|--------------------------------|-------------|------------------------|
    |    NAMESPACE    |              NAME              | TARGET PORT |          URL           |
    |-----------------|--------------------------------|-------------|------------------------|
    | userdataapp-dev | service-userdataapp-backend-np |             | http://127.0.0.1:64874 |
    |-----------------|--------------------------------|-------------|------------------------|
    * Opening service userdataapp-dev/service-userdataapp-backend-np in default browser...
    ! Because you are using a Docker driver on windows, the terminal needs to be open to run it.
    </pre>

It will execute the backend service on default browser. Pick the port from browser.

- Execute Postman and import the Postman collection (`UserDataApp_URLs.postman_collection.json`) from directory `Postman_file/`.
Update the port, picked from browser, in each endpoint's http API requests.

- Each endpoints can be tried in Postman.

   Current user, set through secret in backend application, doesn't have permissions to Select and Insert data.

   So, we will see errors like below in /getUserInfo endpoint:
   <pre> Unable to retrive UserInfo.
    Error Message: mssql: The SELECT permission was denied on the object 'USERINFO', database 'USERDATA', schema 'dbo'.
   </pre>

### Upgrade the application by changing the Helm chart values and performing a Helm upgrade

- Uncomment "select: user2" (at line #35) in values.yaml of backend chart (/userdataapp_backend_chart/values.yaml) and comment "select: user1" (at line #34) then save it.

- Now upgrade the backend application using **`helm upgrade`**:
    <pre> >helm upgrade userdataapp-backend .\userdataapp_backend_chart\ -n userdataapp-dev --set environment.env_dev=true 
    
    Release "userdataapp-backend" has been upgraded. Happy Helming!
    NAME: userdataapp-backend
    LAST DEPLOYED: Mon May  5 00:34:17 2025
    NAMESPACE: userdataapp-dev
    STATUS: deployed
    REVISION: 2
    TEST SUITE: None
    </pre>

- Now try sending request in each specified http requests on Postman. It should work. (You can try /addUser and then /getUserInfo enpoints.)

### Rollback the application to a previous version using Helm rollback

- Now **`Rollback`** the backend application to the **`latest previous revision`** using below command:
    <pre> >helm rollback userdataapp-backend -n userdataapp-dev
    Rollback was a success! Happy Helming!</pre>

    Now try /getUserInfo endpoint again on Postman. You will see the same permission issue again.

- **`helm revision history`** can be seen using below command:
    <pre> >helm history userdataapp-backend -n userdataapp-dev 
    REVISION        UPDATED                         STATUS          CHART                           APP VERSION     DESCRIPTION
    1               Sun May  4 22:56:40 2025        superseded      userdataapp_backend_chart-1.0.0 1.4             Install complete
    2               Mon May  5 00:34:17 2025        superseded      userdataapp_backend_chart-1.0.0 1.4             Upgrade complete
    3               Mon May  5 00:44:08 2025        deployed        userdataapp_backend_chart-1.0.0 1.4             Rollback to 1</pre>

- Now **`Rollback`** the backend application again to **`a previous revision`**:

    Rollback again to go back to the working user replica pod. (2 is the revision number)
    <pre> >helm rollback userdataapp-backend 2 -n userdataapp-dev
    Rollback was a success! Happy Helming!</pre>

    Now try on Postman again. The permission issue must be resolved again.


### `Bonus Question:` Implementation of persistent storage using Kubernetes PersistentVolumes and PersistentVolumeClaims to ensure data persistence across pod restarts

- For persistence volume, I have added following tags in our db helm chart: 

**`.\Helm Charts\userdataapp_db_chart\values.yaml`**
<pre>
pvcInfo:
  pvcName: userdataapp-db-pvc
  storage: ""

volumes:
- name: userdataapp-db-volume
  persistentVolumeClaim:
      claimName: userdataapp-db-pvc

volumeMounts:
- name: userdataapp-db-volume
  mountPath: /var/opt/mssql
</pre>

**`.\Helm Charts\userdataapp_db_chart\templates\persistentVolumeClaim.yaml`**
<pre>
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ .Values.pvcInfo.pvcName }}
  annotations:
    "helm.sh/resource-policy": keep     # this is to keep the persistent volume event after helm uninstall
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Values.pvcInfo.storage | default "1Gi" }}
</pre>

**`.\Helm Charts\userdataapp_db_chart\templates\deployment.yaml`**
<pre>
containers:
  - name: {{ .Values.containerName }}
    -------
    -------
    {{- with .Values.volumeMounts }}
    volumeMounts:
      {{- toYaml . | nindent 12 }}
    {{- end }}

    {{- with .Values.volumes }}
    volumes:
      {{- toYaml . | nindent 12 }}
    {{- end }}
</pre>

#### To check data persistence
- Add some records in database using `UserDataApp_Backend_addUser` http API in Postman (Don't forget to update port number from the service you started using minikube)

#### 1. Rollout Resart DB deployment:

- Get deployment name:
    <pre> >kubectl get deployments -n userdataapp-dev -l app=userdataapp-db 
    NAME             READY   UP-TO-DATE   AVAILABLE   AGE
    userdataapp-db   1/1     1            1           64m</pre>

- Restart db deployment:
    <pre> >kubectl rollout restart deployment userdataapp-db -n userdataapp-dev
    deployment.apps/userdataapp-db restarted</pre>

    (wait for a min or so to let the pod started again)

- Get pod name:
    <pre> >kubectl get pods -n userdataapp-dev -l app=userdataapp-db 
    NAME                              READY   STATUS    RESTARTS      AGE
    userdataapp-db-6557db868b-zxs7g   1/1     Running   2 (95s ago)   2m3s</pre>

- Test and Verify data persistance on Postman on `UserDataApp_Backend_getUserInfo` http API.

#### 2. Uninstall and install db helm chart again:
- Uninstall db helm release (`userdataapp-db`):
    <pre> >helm uninstall userdataapp-db -n userdataapp-dev 
    These resources were kept due to the resource policy:
    [PersistentVolumeClaim] userdataapp-db-pvc


      release "userdataapp-db" uninstalled</pre>

- Install db helm chart again:
    <pre> >helm install userdataapp-db ".\userdataapp_db_chart" -n userdataapp-dev
    NAME: userdataapp-db
    LAST DEPLOYED: Mon May  5 12:28:19 2025
    NAMESPACE: userdataapp-dev
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None</pre>

    (wait for a min or so to let the pod started again)

- Get pod name:
    <pre> >kubectl get pods -n userdataapp-dev -l app=userdataapp-db 
    NAME                              READY   STATUS    RESTARTS   AGE
    userdataapp-db-75955c9759-dzs7l   1/1     Running   0          2m36s</pre>

- Test and Verify data persistance on Postman on `UserDataApp_Backend_getUserInfo` http API.

### Implementation of a namespace strategy for Kubernetes cluster, organizing applications and resources into separate namespaces based on their environments (e.g., development, staging, production)

- Create namespaces:

    - `userdataapp-dev`: namespace for Development (Already created)
    - `userdataapp-uat`: namespace for Staging
    - `userdataapp-prod`: namespace for Production
    <pre> 
    >kubectl create ns userdataapp-uat
    namespace/userdataapp-uat created

    >kubectl create ns userdataapp-prod
    namespace/userdataapp-prod created
    </pre>

#### Environment - Staging
- helm install database resources in Kubernetes cluster:
    <pre> >helm install userdataapp-db ".\userdataapp_db_chart" -n userdataapp-uat
    NAME: userdataapp-db
    LAST DEPLOYED: Mon May  5 13:01:18 2025
    NAMESPACE: userdataapp-uat
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None</pre> 

- helm install backend application resources in Kubernetes cluster:
    <pre> >helm install userdataapp-backend .\userdataapp_backend_chart\ -n userdataapp-uat --set environment.env_uat=true
    NAME: userdataapp-backend
    LAST DEPLOYED: Mon May  5 13:03:47 2025
    NAMESPACE: userdataapp-uat
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None</pre>

- Start backend NodePort service:
    <pre> >minikube service service-userdataapp-backend-np -n userdataapp-uat
    |-----------------|--------------------------------|-------------|---------------------------|
    |    NAMESPACE    |              NAME              | TARGET PORT |            URL            |
    |-----------------|--------------------------------|-------------|---------------------------|
    | userdataapp-uat | service-userdataapp-backend-np |        3031 | http://192.168.49.2:31919 |
    |-----------------|--------------------------------|-------------|---------------------------|
    🏃  Starting tunnel for service service-userdataapp-backend-np.
    |-----------------|--------------------------------|-------------|------------------------|
    |    NAMESPACE    |              NAME              | TARGET PORT |          URL           |
    |-----------------|--------------------------------|-------------|------------------------|
    | userdataapp-uat | service-userdataapp-backend-np |             | http://127.0.0.1:51662 |
    |-----------------|--------------------------------|-------------|------------------------|
    🎉  Opening service userdataapp-uat/service-userdataapp-backend-np in default browser...
    ❗  Because you are using a Docker driver on windows, the terminal needs to be open to run it.</pre>

    This will start the service and run it in default browser. There you will see a response as:

    `Welcome to Kubernetes Advanced Assignment's backend application running on` **`Staging`**

- Test and verify backend application's http APIs on Postman by `changing port number` from this service.

### Environment - Production
- helm install database resources in Kubernetes cluster:
    <pre> >helm install userdataapp-db ".\userdataapp_db_chart" -n userdataapp-prod
    NAME: userdataapp-db
    LAST DEPLOYED: Mon May  5 13:16:10 2025
    NAMESPACE: userdataapp-prod
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None</pre> 

- helm install backend application resources in Kubernetes cluster:
    <pre> >helm install userdataapp-backend .\userdataapp_backend_chart\ -n userdataapp-prod --set environment.env_prod=true
    NAME: userdataapp-backend
    LAST DEPLOYED: Mon May  5 13:17:24 2025
    NAMESPACE: userdataapp-prod
    STATUS: deployed
    REVISION: 1
    TEST SUITE: None</pre>

- Start backend NodePort service:
    <pre> >minikube service service-userdataapp-backend-np -n userdataapp-prod
    |------------------|--------------------------------|-------------|---------------------------|
    |    NAMESPACE     |              NAME              | TARGET PORT |            URL            |
    |------------------|--------------------------------|-------------|---------------------------|
    | userdataapp-prod | service-userdataapp-backend-np |        3031 | http://192.168.49.2:30151 |
    |------------------|--------------------------------|-------------|---------------------------|
    🏃  Starting tunnel for service service-userdataapp-backend-np.
    |------------------|--------------------------------|-------------|------------------------|
    |    NAMESPACE     |              NAME              | TARGET PORT |          URL           |
    |------------------|--------------------------------|-------------|------------------------|
    | userdataapp-prod | service-userdataapp-backend-np |             | http://127.0.0.1:51892 |
    |------------------|--------------------------------|-------------|------------------------|
    🎉  Opening service userdataapp-prod/service-userdataapp-backend-np in default browser...
    ❗  Because you are using a Docker driver on windows, the terminal needs to be open to run it.</pre>

    This will start the service and run it in default browser. There you will see a response as:

    `Welcome to Kubernetes Advanced Assignment's backend application running on` **`Production`**

- Test and verify backend application's http APIs on Postman by `changing port number` from this service.


## Clear everyting from Kubernetes cluster
#### Environment - Development

- Uninstall the helm releases:
    - Uninstall backed application's release:
        <pre> >helm uninstall userdataapp-backend -n userdataapp-dev
        release "userdataapp-backend" uninstalled</pre>

    - Uninstall database release:
        <pre> >helm uninstall userdataapp-db -n userdataapp-dev
        These resources were kept due to the resource policy:
        [PersistentVolumeClaim] userdataapp-db-pvc

          release "userdataapp-db" uninstalled</pre>

- delete namespace (`userdataapp-dev`):
  <pre> >kubectl delete ns userdataapp-dev
    namespace "userdataapp-dev" deleted</pre>

#### Environment - Staging
- Uninstall the helm releases:
    - Uninstall backed application's release:
        <pre> >helm uninstall userdataapp-backend -n userdataapp-uat
        release "userdataapp-backend" uninstalled</pre>

    - Uninstall database release:
        <pre> >helm uninstall userdataapp-db -n userdataapp-uat
        These resources were kept due to the resource policy:
        [PersistentVolumeClaim] userdataapp-db-pvc

          release "userdataapp-db" uninstalled</pre>

- delete namespace (`userdataapp-uat`):
  <pre> >kubectl delete ns userdataapp-uat
    namespace "userdataapp-uat" deleted</pre>

#### Environment - Production
- Uninstall the helm releases:
    - Uninstall backed application's release:
        <pre> >helm uninstall userdataapp-backend -n userdataapp-prod
        release "userdataapp-backend" uninstalled</pre>

    - Uninstall database release:
        <pre> >helm uninstall userdataapp-db -n userdataapp-prod
        These resources were kept due to the resource policy:
        [PersistentVolumeClaim] userdataapp-db-pvc

          release "userdataapp-db" uninstalled</pre>

- delete namespace (`userdataapp-prod`):
  <pre> >kubectl delete ns userdataapp-prod
    namespace "userdataapp-prod" deleted</pre>