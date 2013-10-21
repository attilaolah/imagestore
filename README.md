# Google App Engine image store

`gaeimagestore` implements a very simple image storage service on top of [Google
App Engine][1].

* Images are stored in Blobstore.
* Images are served using the Image Hosting API.

[1]: https://developers.google.com/appengine

The idea is that this app can be deployed in a few minutes and be ready for use
to test other apps. A fair amount of images can easily fit in the free quota.
