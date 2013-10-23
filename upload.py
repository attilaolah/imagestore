#!/usr/bin/env python

import argparse
import os
import re

import requests


SERVER = 'https://gaeimagestore.appspot.com/'
DEFAULT_BATCH_SIZE = 20

RX = re.compile(r'^[0-9a-f]{40}\.jpg', re.I)
RX_UPLOAD = re.compile(r'^https://gaeimagestore.appspot.com/_ah/upload/')


parser = argparse.ArgumentParser()
parser.add_argument('path', help='directory containing the images',
    type=unicode)
parser.add_argument('-c', '--cookie', help='authentication cookie',
    type=bytes, required=True)
parser.add_argument('-b', '--batch', help='batch size (default is 20)',
    type=int)
args = parser.parse_args()


class LoginRequired(ValueError):
    """Raised when login is required."""


def main():
    batch = []
    batch_size = DEFAULT_BATCH_SIZE
    if args.batch is not None:
        batch_size = args.batch
    for img in sorted(os.listdir(args.path)):
        if not RX.match(img):
            continue
        url_path = os.path.join(args.path, img.replace('.jpg', '.url'))
        if os.path.exists(url_path):
            continue
        if len(batch) < batch_size:
            batch.append(img)
            continue
        upload(batch)
        batch = []
    upload(batch)


def upload(batch, path=args.path):
    if not batch:
        return

    response = requests.get(SERVER, allow_redirects=False, headers={
        'Cookie': 'ACSID={}'.format(args.cookie),
    })
    url = response.headers.get('location')
    if not url:
        print 'Location header not found.'
        raise ValueError(response.headers)
    if not RX_UPLOAD.match(url):
        print 'You need to log in.'
        raise LoginRequired(url)

    response = requests.post(url, files=[
        ('file', (img, open(os.path.join(path, img), 'rb')))
    for img in batch])

    for img, url in zip(batch, response.json()):
        if not url:
            continue
        url_path = os.path.join(args.path, img.replace('.jpg', '.url'))
        with open(url_path, 'wb') as fp:
            fp.write(url)
        print img, url


if __name__ == '__main__':
    main()
