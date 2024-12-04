<script lang="ts">
  import {
    NetService,
    PacketTypes,
    type DrawPoint,
    type GameSettingsPacket,
  } from "../../service/net";
  import { onMount } from "svelte";
  let canvas: HTMLCanvasElement;
  let ctx: CanvasRenderingContext2D;
  let colorPicker: HTMLInputElement;
  let lineWidthInput: HTMLInputElement;
  let fillBtn: HTMLButtonElement;

  let isDrawing = false;
  let lastX = 0;
  let lastY = 0;
  let netService: NetService = NetService.getInstance();
  netService.onPacket((packet) => {
    console.log("PACKET ID", packet.id);

    switch (packet.id) {
      case PacketTypes.Coordinates: {
        const data = packet as DrawPoint;
        drawLine(
          data.x1,
          data.y1,
          data.x2,
          data.y2,
          data.color,
          data.lineWidth
        );
        break;
      }

      case PacketTypes.GameSettings: {
        const data = packet as GameSettingsPacket;

        for (const point of data.coordinates) {
          drawLine(
            point.x1,
            point.y1,
            point.x2,
            point.y2,
            point.color,
            point.lineWidth
          );
        }
        break;
      }
    }

    // const  const data cordPacket = packet as DrawPoint;
    // const data = cordPacket;
  });

  onMount(() => {
    // Ensure canvas and context are available after component mounts
    if (canvas) {
      ctx = canvas.getContext("2d") as CanvasRenderingContext2D;
      if (!ctx) {
        console.error("Failed to get canvas context");
        return;
      }
    } else {
      console.error("Canvas element not found");
    }

    // Handle incoming messages from the server

    // Fill canvas with a color when button is clicked
    if (fillBtn) {
      fillBtn.addEventListener("click", () => {
        ctx.fillStyle = "blue"; // Set fill color to blue
        ctx.fillRect(0, 0, canvas.width, canvas.height); // Fill the entire canvas
      });
    }

    // Add event listeners for mouse events
    if (canvas) {
      canvas.addEventListener("mousedown", (e) => {
        isDrawing = true;
        lastX = e.offsetX;
        lastY = e.offsetY;

        const message = {
          x1: lastX,
          y1: lastY,
          x2: lastX,
          y2: lastY,
          color: colorPicker.value,
          lineWidth: lineWidthInput.value,
        };

        let packet: DrawPoint = {
          id: PacketTypes.Coordinates,
          ...message,
        };
        netService.sendPacket(packet);
      });

      canvas.addEventListener("mouseup", () => {
        isDrawing = false;
        ctx.beginPath();
      });

      canvas.addEventListener("mousemove", draw);
    }
  });

  // Drawing function to send data to the server and draw locally
  function draw(e: MouseEvent) {
    if (!isDrawing) return;

    ctx.lineWidth = parseFloat(lineWidthInput.value);
    ctx.lineCap = "round";
    ctx.strokeStyle = colorPicker.value;

    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    const message = {
      x1: lastX,
      y1: lastY,
      x2: x,
      y2: y,
      color: colorPicker.value,
      lineWidth: lineWidthInput.value,
    };

    let packet: DrawPoint = {
      id: PacketTypes.Coordinates,
      ...message,
    };
    netService.sendPacket(packet);

    drawLine(lastX, lastY, x, y, colorPicker.value, lineWidthInput.value);

    lastX = x;
    lastY = y;
  }

  // Helper function to draw the line on canvas
  function drawLine(
    x1: number,
    y1: number,
    x2: number,
    y2: number,
    color: string,
    lineWidth: string
  ) {
    ctx.lineWidth = parseFloat(lineWidth);
    ctx.lineCap = "round";
    ctx.strokeStyle = color;

    ctx.beginPath();
    ctx.moveTo(x1, y1);
    ctx.lineTo(x2, y2);
    ctx.stroke();
  }
</script>

<div class="h-screen w-screen p-20">
  <div class="flex flex-col w-full h-full justify-center items-center">
    <canvas
      bind:this={canvas}
      id="drawingCanvas"
      width="800"
      height="600"
      class="border-black border-[1px] cursor-crosshair flex justify-center flex-1"
    ></canvas>
    <br />
    <div class="w-[25%] flex flex-col gap-4 items-center justify-center">
      <div class="flex gap-4 items-center w-full justify-center">
        <p class="text-sm w-[30%]">Color picker:</p>
        <input
          bind:this={colorPicker}
          type="color"
          id="colorPicker"
          value="#000000"
          class="flex-1"
        />
      </div>

      <div class="flex gap-4 items-center w-full justify-center">
        <p class="text-sm w-[30%] justify-center">Line Width:</p>
        <input
          bind:this={lineWidthInput}
          type="number"
          id="lineWidth"
          value="2"
          min="1"
          max="100"
          placeholder="line width"
          class="flex-1 text-xs py-2 px-2 border-2 rounded-md m-0 focus:outline-0"
        />
      </div>
    </div>
  </div>
</div>
