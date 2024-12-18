<script lang="ts">
  import {
    NetService,
    PacketTypes,
    type DrawPoint,
  } from "../../service/net";
  import { onDestroy, onMount } from "svelte";

  let canvas: HTMLCanvasElement;
  let ctx: CanvasRenderingContext2D;
  let colorPicker: HTMLInputElement;
  let lineWidthInput: HTMLInputElement;
  let fillBtn: HTMLButtonElement;
  let isDrawing = false;
  let lastX = 0;
  let lastY = 0;
  export let pointerEvents: boolean;
  export let netService: NetService;
  export let exposeMethods;

  $: if (exposeMethods) {
    exposeMethods({ drawLine,resetCanvas });
  }

  // const gameCode: Writable<string> = writable("");

  onDestroy(() => {
    console.log("Destroying player view...");
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
        console.log("Sending packet");
        netService.sendPacket(packet);
      });

      canvas.addEventListener("mouseup", () => {
        isDrawing = false;
        ctx.beginPath();
      });

      canvas.addEventListener("mousemove", draw);
    }

    console.log("View got mounted...");
  });

  function resetCanvas() {
    if (ctx && canvas) {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
    }
  }

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
  export function drawLine(
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

<div
  class={`h-screen w-screen p-20 ${pointerEvents ? "pointer-events-auto" : "pointer-events-none"}`}
>
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
