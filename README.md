# Face-Recognition-Golang and Python
Face recognition using Go and Python

The images are taken as input by the GO microservice, and is sent to the Python microservice. The Python microservice, then processes the distance between the two images using `face_recognition`
package and returns the distance in json format. The output from the python microservice is collected by go microservice and the output is displayed.

The images taken as input are stored in the GCP.
