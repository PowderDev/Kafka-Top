import {
  Chart as ChartJS,
  Legend,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  TimeScale,
} from "chart.js"
import { Line } from "react-chartjs-2"
import "chartjs-adapter-date-fns"
import { useEffect, useRef, useState } from "react"
import "./App.css"
import { randomHexColor } from "./helpers"

ChartJS.register(Legend, CategoryScale, LinearScale, PointElement, LineElement, TimeScale)

type Dataset = {
  label: string
  data: { x: number; y: number }[]
  backgroundColor?: string
  borderColor?: string
}

type wsData = {
  type: string
  data: { ts: number; value: number[] }
}

function App() {
  const socket = useRef<WebSocket | null>(null)
  const [datasets, setDatasets] = useState<Dataset[]>([])
  const [wsData, setWsData] = useState<wsData>({} as wsData)

  useEffect(() => {
    socket.current = new WebSocket("ws://backend.kafka-top.local/ws")

    socket.current.onopen = () => {
      console.log("Connected to the server")
    }

    socket.current.onmessage = (event) => {
      setWsData(JSON.parse(event.data))
    }

    socket.current.onclose = () => {
      console.log("Disconnected from the server")
    }

    socket.current.onerror = (error) => {
      console.error(error)
    }

    return () => {
      socket.current?.close()
    }
  }, [])

  useEffect(() => {
    if (wsData.type === "cpu_percent") {
      const xyArray = wsData.data.value.map((value) => {
        return { x: wsData.data.ts, y: value }
      })

      const newDatasets = [...datasets]

      for (let i = 0; i < wsData.data.value.length; i++) {
        if (newDatasets.length <= i) {
          const color = randomHexColor()
          newDatasets.push({
            label: `Core ${i + 1}`,
            data: [],
            backgroundColor: color,
            borderColor: color,
          })
        }

        newDatasets[i].data.push(xyArray[i])
        newDatasets[i].data = newDatasets[i].data.slice(-10)
      }

      setDatasets(newDatasets)
      setWsData({} as wsData)
    }
  }, [wsData])

  return (
    <div id="root">
      <h1>Top Dashboard</h1>

      <div style={{ width: "600px", margin: "0 auto" }}>
        <Line
          data={{ datasets }}
          options={{
            scales: {
              x: {
                type: "time",
                time: {
                  minUnit: "second",
                  displayFormats: {
                    second: "h:mm:ss",
                  },
                },
                ticks: {
                  source: "data",
                },
              },
            },
          }}
        />
      </div>
    </div>
  )
}

export default App
