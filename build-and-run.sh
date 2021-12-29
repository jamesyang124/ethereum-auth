#!/bin/bash

#!/bin/bash

REGISTRY="585958861587.dkr.ecr.us-west-2.amazonaws.com"
PROJECT="identity/ethereum-auth"
BASE="${REGISTRY}/${PROJECT}"

if [ $(git rev-parse HEAD) = $(git rev-parse origin/master) ]; then
  TAG="1.0.$(date -u '+%y%m%d%H%M')-$(git describe --long --dirty --abbrev=10 --tags --always)"
elif [ $(git rev-parse HEAD) = $(git rev-parse origin/prod) ]; then
  TAG="1.1.$(date -u '+%y%m%d%H%M')-$(git describe --long --dirty --abbrev=10 --tags --always)"
else
  TAG="1.0.$(date -u '+%y%m%d%H%M')-$(git describe --long --dirty --abbrev=10 --tags --always)-QA"
fi

# you must commit changes first
git diff --cached --exit-code || exit 1
git diff --exit-code || exit 2

docker build --squash --build-arg "SGNGX_CONFIG_VERSION=${TAG}" -t "${BASE}:${TAG}" .  \
|| docker build --build-arg "SGNGX_CONFIG_VERSION=${TAG}" -t "${BASE}:${TAG}" .  \

IS_AWS_V2=$(aws --version | grep "aws-cli/2")

if [[ $IS_AWS_V2 == "" ]]; then
  $(aws ecr get-login --no-include-email --region us-west-2 --registry-ids 585958861587)
else
  R=us-west-2; A=585958861587; aws ecr  get-login-password  --region "${R}" |  docker login -u AWS --password-stdin "${A}.dkr.ecr.${R}.amazonaws.com"
fi

[ $? -eq 0 ] && docker push "${BASE}:${TAG}"

if [ $? -eq 0 ]; then

echo
echo "${BASE}:${TAG}"
echo
echo "##vso[task.setvariable variable=BUILT_IMAGE_REASON]${BUILD_SOURCEVERSIONMESSAGE}"
echo "##vso[task.setvariable variable=BUILT_IMAGE]/${PROJECT}:${TAG}" # feedback to VSTS
echo "#############################################"
echo "##"
echo "## IMAGE:  /${PROJECT}:${TAG}"
echo "##"
echo "## WHY: ${BUILD_SOURCEVERSIONMESSAGE}"
echo "##"
echo "## Pushed to: ${REGISTRY}"
[ ! -z "${REGISTRY2}" ] && echo "## Pushed to: ${REGISTRY2}"
echo "##"
echo "#############################################"
echo
echo
echo
exit 0

else
  echo
  echo "docker push back to dev regsitry failed!"
  echo "do you login dev registry before?"
  echo
  exit 1
fi