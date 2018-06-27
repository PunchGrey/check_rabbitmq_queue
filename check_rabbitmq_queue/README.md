# check_rabbitmq_queue
check_rabbitmq_queue
====================
Плагин написан на go, проверяет по определенным параметрам количество сообщений и количество неподтвержденных сообщений в очереди. Проверка работает только, если RabbitMQ предоставляет доступ через web. 

Проверять можно определенную очередь (-q name_queue), или несколько очередей (-q name1_queue,name2_queue), или все очереди (-q ALL). Но независимо от варианта проверки, параметры определяющие допустимые значения задаются для всех очередей одни, их нельзя устанавливать для каждой очереди. 

Применение
----------
    check_queue -h 127.0.0.1 -p 15672 -u admin -pw 1234 -q celery -em 28 -wm 15 -emu 5 -wmu 2
        опции:
            -h хост, на котором расположен RabbitMQ
                значение по умолчанию 127.0.0.1
            -p порт, через который доступен RabbitMQ (rabbitmq_management)
                значение по умолчанию 15672
            -u пользователь
                значение по умолчанию admin
            -pw пароль
                значение по умолчанию 1234
            -q имя очереди (-q name_queue), или несколько очередей (-q name1_queue,name2_queue), или все очереди (-q ALL)
                значение по умолчанию celery
            -em если больше этого значения сообщений в очереди, то возвращает ERROR
                значение по умолчанию 28
            -wm если больше этого значения сообщений в очереди, то возвращает WARNING
                значение по умолчанию 15
            -emu если больше этого значения неподтвержденных сообщений в очереди, то возвращает ERROR
                значение по умолчанию 5
            -wmu если больше этого значения неподтвержденных сообщений в очереди, то возвращает WARNING
                значение по умолчанию 2
        вывод:
            name queue celery: 28 messages,  0 MessagesUnacknowledged
            exit status 1        

Прописываем в commands.cfg

    define command {
        command_name    check_queue_celery
        command_line    $USER1$/check_nrpe -H $HOSTADDRESS$ -c check_queue -a $ARG1$ $ARG2$ $ARG3$ $ARG4$ $ARG5$ $ARG6$ $ARG7$ $ARG8$ $ARG9$
    }

Прописываем в service.cfg
    define service {
        check_command                  check_queue_celery!testhost.local!127.0.0.1!15672!admin!1234!celery!28!15!5!2
        contact_groups                 admins
        host_name                      testhost.local
        service_description            RABBITMQ::check quaue::test.host.local
        servicegroups                  core
        use                            local-service
    }

На хосте testhost.local (откуда будет запускаться проверка)
в файле /etc/nrpe.d/check_queue.cfg
    command[check_queue_celery]=/usr/lib64/nagios/plugins/check_queue -h $ARG1$ -p $ARG2$ -u $ARG3$ -pw $ARG4$ -q $ARG5$ -em $ARG6$ -wm $ARG7$ -emu $ARG8$ -wmu $ARG9$
/usr/lib64/nagios/plugins/check_queue - бинарник должен быть на testhost.local