FROM public.ecr.aws/lambda/python:3.11

COPY ./docker/requirements.txt /var/task
RUN pip install -r /var/task/requirements.txt
COPY ../src /var/task

CMD [ "lambda_function.lambda_handler" ]