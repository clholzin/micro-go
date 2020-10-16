

### flash it

```
cp ./s110_nrf51_8.0.0/s110_nrf51_8.0.0_softdevice.hex /Volume/MICROBIT/
```

### run it

```
tinygo flash -target=microbit-s110v8 main.go

```

saw different flag programmer
```
tinygo flash -target=microbit-s110v8 -programmer=cmsis-dap
```
