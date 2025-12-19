"use client"

import { useEffect } from "react"
import { toast } from "sonner"
import { useWebSocket } from "@/hooks/use-websocket"

export function RealTimeNotifications() {
  // Assuming the backend is running on localhost:8080
  // In production, this should be an environment variable
  const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/ws"
  
  const { lastMessage, isConnected } = useWebSocket(WS_URL)

  useEffect(() => {
    if (isConnected) {
      // toast.success("Connected to real-time updates")
    }
  }, [isConnected])

  useEffect(() => {
    if (lastMessage) {
      try {
        // Try to parse as JSON if possible, or just show string
        const data = JSON.parse(lastMessage)
        toast.info("New Notification", {
          description: data.message || JSON.stringify(data),
        })
      } catch (e) {
        toast.info("New Notification", {
          description: lastMessage,
        })
      }
    }
  }, [lastMessage])

  return null
}
