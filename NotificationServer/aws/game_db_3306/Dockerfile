# 베이스 이미지로 Ubuntu 22.04 사용
FROM ubuntu:22.04

# 작업 디렉토리 설정
WORKDIR /app

# 호스트의 game_server 파일을 복사
COPY game_server .

# 8080 포트 노출
EXPOSE 8080

# 컨테이너 실행 시 game_server 실행
CMD ["./game_server"]