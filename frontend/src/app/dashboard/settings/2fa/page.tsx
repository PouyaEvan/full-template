"use client"

import { useState } from "react"
import axios from "axios"
import { toast } from "sonner"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { zodResolver } from "@hookform/resolvers/zod"

const twoFactorSchema = z.object({
  code: z.string().length(6, {
    message: "Code must be 6 digits.",
  }),
})

export default function TwoFactorSetupPage() {
  const [qrCode, setQrCode] = useState<string | null>(null)
  const [backupCodes, setBackupCodes] = useState<string[]>([])
  const [isEnabled, setIsEnabled] = useState(false)

  const form = useForm<z.infer<typeof twoFactorSchema>>({
    resolver: zodResolver(twoFactorSchema),
    defaultValues: { code: "" },
  })

  const setup2FA = async () => {
    try {
      const token = localStorage.getItem("token")
      const res = await axios.post("http://localhost:8080/api/auth/2fa/setup", {}, {
        headers: { Authorization: `Bearer ${token}` },
        responseType: 'blob' // Important for image
      })
      
      const url = URL.createObjectURL(res.data)
      setQrCode(url)
    } catch (error) {
      toast.error("Failed to setup 2FA")
    }
  }

  const enable2FA = async (values: z.infer<typeof twoFactorSchema>) => {
    try {
      const token = localStorage.getItem("token")
      const res = await axios.post("http://localhost:8080/api/auth/2fa/enable", 
        { code: values.code },
        { headers: { Authorization: `Bearer ${token}` } }
      )
      
      setBackupCodes(res.data.backup_codes)
      setIsEnabled(true)
      setQrCode(null)
      toast.success("2FA Enabled Successfully")
    } catch (error) {
      toast.error("Invalid Code")
    }
  }

  return (
    <div className="p-8 max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-6">Two-Factor Authentication</h1>
      
      {!isEnabled && !qrCode && (
        <button 
          onClick={setup2FA}
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Setup 2FA
        </button>
      )}

      {qrCode && !isEnabled && (
        <div className="space-y-6">
          <div className="border p-4 rounded bg-white inline-block">
            <img src={qrCode} alt="2FA QR Code" />
          </div>
          <p>Scan this QR code with your authenticator app (Google Authenticator, Authy, etc.)</p>
          
          <form onSubmit={form.handleSubmit(enable2FA)} className="space-y-4 max-w-xs">
            <input
              {...form.register("code")}
              className="w-full border p-2 rounded"
              placeholder="Enter 6-digit code"
            />
            <button 
              type="submit"
              className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 w-full"
            >
              Verify & Enable
            </button>
          </form>
        </div>
      )}

      {isEnabled && (
        <div className="space-y-4">
          <div className="bg-green-100 text-green-800 p-4 rounded">
            2FA is currently enabled on your account.
          </div>
          
          {backupCodes.length > 0 && (
            <div className="border p-4 rounded">
              <h3 className="font-bold mb-2">Backup Codes</h3>
              <p className="text-sm text-gray-600 mb-4">Save these codes in a safe place. You can use them to login if you lose access to your authenticator app.</p>
              <div className="grid grid-cols-2 gap-2 font-mono bg-gray-100 p-4 rounded">
                {backupCodes.map((code, i) => (
                  <div key={i}>{code}</div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
