# producer
    > send task to the broker(message queue system)
# broker
    > the message queue system, FIFO
# worker
    > get message from broker, and consume this message
# backend
    > restore the message handle result

## TODO
- More ways to config backend and broker(ENV, config yml e.g.)
- Add support for delay task
- Add support for periodic task
