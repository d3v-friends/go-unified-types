# go-unified-types

* 모든 proto 정의에서 공용으로 사용하는 타입
    * mongo-driver, gqlgen, gorm 에서 모두 사용 가능한 타입
* utypes.proto 파일은 protoc 이 설치된 include 디렉터리에 복사해둬야 code generate 가 된다.
* github.com/d3v-friends/mango 의 쿼리관련 기능이 구현되어 있다.
* github.com/d3v-friends/go-tools 를 사용한다.

~~~
	$PROTOC_PATH/include/d3v-friends/go-unified-types/utypes.proto
~~~