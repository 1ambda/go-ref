#!/usr/bin/env bash

TAG="[generate-gomock.sh]"

MOCK_PREFIX="mock"

GO_FILES=$(ls ${PKG_DIR}/*.go | xargs -n 1 basename)

for GO_FILE in ${GO_FILES[@]}
do
    PKG=$(basename ${PKG_DIR})
    MOCK_FILE=${PKG_DIR}/${MOCK_PREFIX}_${GO_FILE}
    mockgen -package ${PKG} -source ${PKG_DIR}/${GO_FILE} -destination ${MOCK_FILE}

    # remove if the generated file doesn't include `struct` (useless mock files)
    if cat ${MOCK_FILE} | grep -q --line-buffered struct;
    then
        echo -e "Generated: ${MOCK_FILE}"
    else
        rm ${MOCK_FILE}
    fi

done
