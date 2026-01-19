# Introducción

## Objetivo del sistema
Proveer una base de observabilidad para eventos y telemetría de drones, desacoplando el edge del backend mediante un bus MQTT y exponiendo métricas Prometheus para análisis y alertado en Grafana.

## Problema que resuelve
- Centraliza eventos críticos y señales operativas sin acoplar el productor (drone/edge) al consumidor (backend).
- Permite inspección y alertas sin instrumentación propietaria.
- Establece un punto único para medir salud y comportamiento del sistema.

## Público objetivo
- Equipos de ingeniería y operaciones que necesitan observabilidad básica y extensible.
- Proyectos académicos/TFG que requieren una arquitectura realista y justificable.
- Integradores IoT que buscan una base reproducible en Docker.

## Qué NO es este sistema
- No es un sistema de control de vuelo ni un autopiloto.
- No es una plataforma de analítica avanzada ni un data lake.
- No reemplaza un SOC completo ni un SIEM; solo expone señales base.
