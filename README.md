# HEKA-UI 

## BUILD STEPS

1. Clone this repository to your local workspace.
2. Go inside the cloned directory.
3. Generate the docker image.
   * Run `make docker`
4. See the list of generated docker image.
   * Run `docker images infobloxcto/heka-ui`
5. Create a portable docker image file from the generated docker image. 
   * Run `docker save infobloxcto/heka-ui:$(git describe --always --tags) -o heka-ui.docker.img`

**RESULT:** Generated a portable docker image file _heka-ui.docker.img_ in current directory.

## DEPLOYMENT STEPS 

1. Enable `support_access` in the GRID, where the deployment will take place.
2. Copy the portable docker image file to Grid Master.
   * Run `scp heka-ui.docker.img root@<GM_IP>:/root/`
3. SSH in to the GRID Master and go to the /root directory.
   * Run `ssh root@<GM_IP>`
   * Run `cd /root`
4. Load the docker image from portable docker image file, which we transferred here in STEP-2. 
   * Run `docker load -i heka-ui.docker.img`
5. See the list of heka-ui docker images available in NIOS.
   * Run `docker images infobloxcto/heka-ui`
   * Make a note of TAG field values for the image, it will be used in next step.
   * If there are multiple images, select the latest one or the one which needs to be tested.
6. Start the Heka-UI container. Replace the `<TAG>` with TAG value copied in last STEP.
   * Run `docker run -d --net host --name heka_ui --restart=always -v /rw/noa/onprem.d/:/etc/onprem.d/ -v /rw/heka:/etc/storage/hekaui infobloxcto/heka-ui:<TAG>`
7. Check the logs of the Heka-UI container.
   * Run `docker logs -f heka_ui`

**RESULT:** Open any browser application, go to `https://<GM_IP>/bloxconnect` to see Heka-UI login page.

## DEPLOYMENT CLEANUP STEPS

1. Stop and remove running Heka-UI container.
   * Run `docker stop heka_ui`
   * Run `docker rm heka_ui`
2. Remove HTTP injection which happeded during last Heka-UI container run.
   * Run `curl -X POST -H "Content-type: application/json" -H "Accept: application/json" -d '{"name":"remove_http_redirection","args":[]}' "http://0.0.0.0:999/nios/api/v1.0/nios/1/exec"`
3. Wait 1-2 minutes for http server restart.
4. Remove existing docker images and portable docker image files.
   * Run `docker rmi -f $(docker images infobloxcto/heka-ui -q)`
   * Run `rm /root/heka-ui.docker.img`

**RESULT:** This will clean up NIOS from Heka-UI deployment. You can follow **DEPLOYMENT STEPS** section again, to make a new deployment. 
