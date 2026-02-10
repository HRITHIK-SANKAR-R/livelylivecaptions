import asyncio
import websockets
import torch
import numpy as np

# Load the Silero VAD model
model, utils = torch.hub.load(repo_or_dir='snakers4/silero-vad',
                              model='silero_vad',
                              force_reload=True)

(get_speech_timestamps,
 save_audio,
 read_audio,
 VADIterator,
 collect_chunks) = utils

async def handler(websocket):
    print(f"Client connected: {websocket.remote_address}")
    vad_iterator = VADIterator(model)
    audio_buffer = bytearray()
    sample_rate_in = 44100
    sample_rate_out = 16000
    
    try:
        async for message in websocket:
            audio_buffer.extend(message)

            # Process in chunks of a certain size (e.g., 1536 samples for 16kHz)
            # This is a common chunk size for VAD models
            chunk_size = 1536 * (sample_rate_in // sample_rate_out) * 2 # 2 bytes per sample
            
            while len(audio_buffer) >= chunk_size:
                chunk = audio_buffer[:chunk_size]
                audio_buffer = audio_buffer[chunk_size:]

                # The audio data from Go is raw S16LE, so we need to convert it to a tensor
                audio_int16 = np.frombuffer(chunk, np.int16)
                audio_float32 = audio_int16.astype(np.float32) / 32768.0
                
                # Resample the audio
                resample_factor = sample_rate_in // sample_rate_out
                resampled_audio = audio_float32[::resample_factor]

                speech_dict = vad_iterator(torch.from_numpy(resampled_audio), return_seconds=True)
                if speech_dict:
                    print(f"Speech detected at {speech_dict['start']:.2f}s")

            # For now, we'll still send back a dummy caption to keep the client happy
            await websocket.send("This is a dummy caption from the Python server.")

    except websockets.exceptions.ConnectionClosed:
        print(f"Client disconnected: {websocket.remote_address}")
        vad_iterator.reset_states()
    except Exception as e:
        print(f"An error occurred: {e}")

async def main():
    async with websockets.serve(handler, "localhost", 8080):
        print("WebSocket server started at ws://localhost:8080")
        print("Silero VAD model loaded.")
        await asyncio.Future()  # run forever

if __name__ == "__main__":
    # This is needed for Silero VAD to work correctly
    torch.set_num_threads(1)
    asyncio.run(main())