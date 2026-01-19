FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
     python3 \
     python3-venv \
     python3-pip \
     ca-certificates \
  && rm -rf /var/lib/apt/lists/*

# Crear virtualenv explícito (requerido en Ubuntu 24.04+)
RUN python3 -m venv $VIRTUAL_ENV

WORKDIR /app

# Instalar dependencias Python dentro del venv
RUN pip install --no-cache-dir paho-mqtt

COPY mavlink_to_mqtt /app

CMD ["python", "main.py"]
