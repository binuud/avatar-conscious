## Yolo

```
docker run -it -v ./aiData:/src ultralytics/ultralytics:latest-cpu yolo predict model=yolo26n.pt source='/src/street.png'
```