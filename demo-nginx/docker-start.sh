
conf_dir=`pwd`
web_page_dir=`pwd`"/web"
# 目录名为docker容器名
container_name=$(basename `pwd`)

echo "---- start $container_name ... "
echo "---- conf_dir:$conf_dir"
echo "---- web_page_dir:$web_page_dir"

container_id=`docker ps -aq --filter name=${container_name}`
if [ "${container_id}" != "" ]; then
    echo "---- ${container_name} is exists, id:${container_id}"
    # 查看进程是否存在
    status=`docker inspect --format '{{.State.Status}}' ${container_name}`
    echo "---- container:${container_name} status is ${status}"
    if [ "${status}" == "running" ]; then
        echo "---- docker stop $container_name:"
        docker stop ${container_name}
    fi

    echo "---- docker rm $container_name:"
    docker rm ${container_name}
fi

echo "---- docker run $container_name:"
docker run -d -p 80:80 -p 443:443 --privileged\
 -v $web_page_dir:/data/web \
 -v $conf_dir/nginx/conf/conf.d:/usr/local/openresty/nginx/conf/conf.d \
 -v $conf_dir/nginx/conf/nginx.conf:/usr/local/openresty/nginx/conf/nginx.conf \
 -v $conf_dir/nginx/tmp/:/usr/local/openresty/nginx/tmp \
 -v $conf_dir/nginx/logs/:/usr/local/openresty/nginx/logs \
 --name $container_name openresty/openresty:latest

echo "---- SUCCESS!"
