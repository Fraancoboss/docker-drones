# Roadmap técnico

## Evoluciones previstas
- Endurecimiento de seguridad en MQTT y observabilidad.
- Escalado horizontal del backend.
- Soporte multi-dron con aislamiento por topic.

## Integración MAVLink/PX4 (futuro)
- Mapeo de mensajes MAVLink a telemetría/eventos.
- Simulación reproducible sin alterar el core.
- Validación de esquemas antes de exponer métricas.

## Escalado multi-dron
- Namespacing de topics por identificador de drone.
- Shared subscriptions para balancear backends.
- Cuotas de publicación por dispositivo.

## Seguridad avanzada
- TLS mutuo en MQTT.
- Rotación de credenciales y gestión de identidad por drone.
- Firmas de mensajes y detección de replay.

## Gemelo digital / simulación
- Simulación de flotas para pruebas de carga.
- Integración con escenarios de seguridad controlados.
- Métricas sintéticas para validar alertas.

## Estado actual vs evolución
- Estado actual: base estable y mínima.
- Evolución futura: capacidades avanzadas sin romper la arquitectura base.
