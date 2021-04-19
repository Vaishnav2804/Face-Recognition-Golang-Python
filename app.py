from flask import Flask, request
import json

# Face Detection Dependencies
import face_recognition
import matplotlib.pyplot as plt
from matplotlib.patches import Rectangle
from matplotlib.patches import Circle
import numpy as np
import cv2

# GCP Dependencies
from google.cloud import storage

app = Flask(__name__)


def download_blob(bucket_name, source_blob_name1, source_blob_name2):
    """Downloads a blob from the bucket."""
    # bucket_name = "your-bucket-name"
    # source_blob_name = "storage-object-name"
    # destination_file_name = "local/path/to/file"
    
    storage_client = storage.Client.from_service_account_json(
        'face-recognition-311111-b6dce9fd7ca9.json')

    bucket = storage_client.bucket(bucket_name)

    # Construct a client side representation of a blob.
    # Note `Bucket.blob` differs from `Bucket.get_blob` as it doesn't retrieve
    # any content from Google Cloud Storage. As we don't need additional data,
    # using `Bucket.blob` is preferred here.
    blob = bucket.blob(source_blob_name1)
    blob.download_to_filename(source_blob_name1)
    print(
        "Blob {} downloaded to {}.".format(
            source_blob_name1, source_blob_name1
        )
    )
    blob = bucket.blob(source_blob_name2)
    blob.download_to_filename(source_blob_name1)
    print(
        "Blob {} downloaded to {}.".format(
            source_blob_name2, source_blob_name2
        )
    )


@app.route('/')
def index():

    filename1 = request.args.get('pic1')
    filename2 = request.args.get('pic2')

    print("Filename 1 ")
    print(filename1)
    print("Filename 2 ")
    print(filename2)

    download_blob("fr_vaishnav", filename1, filename2)

    image = cv2.imread(filename1)
    pic1 = cv2.cvtColor(image,cv2.COLOR_BGR2RGB)
    imagePic1=pic1
    #face location 

    face_locations=face_recognition.face_locations(imagePic1)

    number_of_faces = len(face_locations)
    print("Known Image")
    print("Found {} face in the input image".format(number_of_faces))

    # plt.imshow(imageSarvesh)
    # ax=plt.gca()

    for face_location in face_locations:
        top,right,bottom,left=face_location
        x,y,w,h = left,top,right,bottom
        print("A face is at pixel location Top:{}, Left:{}, Bottom:{}, Right:{}".format(x,y,w,h))

        # rect= Rectangle((x,y),w-x,h-y,fill=False,color='red')
        # ax.add_patch(rect)
    # plt.show()

    Pic1_encoding=face_recognition.face_encodings(pic1)[0]
    known_face_encodings = [
                        Pic1_encoding,
    ]

    image = cv2.imread(filename2)
    unknown_image = cv2.cvtColor(image,cv2.COLOR_BGR2RGB)
    imagePic2=unknown_image

    face_locations=face_recognition.face_locations(unknown_image)

    number_of_faces = len(face_locations)
    print("Unknown Image")
    print("Found {} face in the input image".format(number_of_faces))

    # plt.imshow(unknown_image)
    # ax=plt.gca()

    for face_location in face_locations:
        top,right,bottom,left=face_location
        x,y,w,h = left,top,right,bottom
        print("A face is at pixel location Top:{}, Left:{}, Bottom:{}, Right:{}".format(x,y,w,h))

    # rect= Rectangle((x,y),w-x,h-y,fill=False,color='red')
    # ax.add_patch(rect)

    unknown_face_encodings=face_recognition.face_encodings(unknown_image)

    from scipy.spatial import distance

    resultJSON = {}

    for unknown_face_encoding in unknown_face_encodings:
        results=[]
        for known_face_encoding in known_face_encodings:
            d=face_recognition.face_distance(known_face_encodings, unknown_face_encoding)
            print("Distance is:",d)
            resultJSON.update({"Distance": d[0]})

    return resultJSON

if __name__ == '__main__':
    app.run(debug=True)