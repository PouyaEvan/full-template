"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useState } from "react"
import axios from "axios"
import { toast } from "sonner"

const formSchema = z.object({
  phone: z.string().min(10, {
    message: "Phone number must be at least 10 characters.",
  }),
})

const otpSchema = z.object({
  code: z.string().length(6, {
    message: "OTP must be 6 digits.",
  }),
})

export default function LoginPage() {
  const [step, setStep] = useState<"phone" | "otp">("phone")
  const [phone, setPhone] = useState("")

  const phoneForm = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      phone: "",
    },
  })

  const otpForm = useForm<z.infer<typeof otpSchema>>({
    resolver: zodResolver(otpSchema),
    defaultValues: {
      code: "",
    },
  })

  async function onPhoneSubmit(values: z.infer<typeof formSchema>) {
    try {
      await axios.post("http://localhost:8080/api/auth/otp/send", { phone: values.phone })
      setPhone(values.phone)
      setStep("otp")
      toast.success("OTP sent successfully")
    } catch (error) {
      toast.error("Failed to send OTP")
    }
  }

  async function onOtpSubmit(values: z.infer<typeof otpSchema>) {
    try {
      const res = await axios.post("http://localhost:8080/api/auth/otp/verify", { phone, code: values.code })
      localStorage.setItem("token", res.data.token)
      toast.success("Login successful")
      // Redirect to dashboard
    } catch (error) {
      toast.error("Invalid OTP")
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 dark:bg-zinc-900">
      <div className="w-full max-w-md space-y-8 rounded-lg border bg-white p-6 shadow-lg dark:border-zinc-800 dark:bg-zinc-950">
        <div className="text-center">
          <h2 className="text-2xl font-bold tracking-tight">
            {step === "phone" ? "Sign in with Phone" : "Enter OTP"}
          </h2>
          <p className="text-sm text-muted-foreground">
            {step === "phone"
              ? "Enter your mobile number to receive a code"
              : `Code sent to ${phone}`}
          </p>
        </div>

        {step === "phone" ? (
          <form onSubmit={phoneForm.handleSubmit(onPhoneSubmit)} className="space-y-4">
            <div>
              <input
                {...phoneForm.register("phone")}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                placeholder="09123456789"
              />
              {phoneForm.formState.errors.phone && (
                <p className="text-sm text-red-500 mt-1">{phoneForm.formState.errors.phone.message}</p>
              )}
            </div>
            <button
              type="submit"
              className="inline-flex h-10 w-full items-center justify-center whitespace-nowrap rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
            >
              Send Code
            </button>
          </form>
        ) : (
          <form onSubmit={otpForm.handleSubmit(onOtpSubmit)} className="space-y-4">
            <div>
              <input
                {...otpForm.register("code")}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                placeholder="123456"
              />
              {otpForm.formState.errors.code && (
                <p className="text-sm text-red-500 mt-1">{otpForm.formState.errors.code.message}</p>
              )}
            </div>
            <button
              type="submit"
              className="inline-flex h-10 w-full items-center justify-center whitespace-nowrap rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50"
            >
              Verify & Login
            </button>
          </form>
        )}
      </div>
    </div>
  )
}
