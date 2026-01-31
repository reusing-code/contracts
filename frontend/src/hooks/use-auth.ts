import { useState, useCallback } from "react"
import { getToken, setToken, clearToken } from "@/lib/api"
import { login as apiLogin, register as apiRegister } from "@/lib/auth-repository"
import type { AuthUser, LoginData, RegisterData } from "@/types/auth"

const USER_KEY = "user"

function loadUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function useAuth() {
  const [user, setUser] = useState<AuthUser | null>(loadUser)
  const [token, setTokenState] = useState<string | null>(getToken)

  const login = useCallback(async (data: LoginData) => {
    const res = await apiLogin(data)
    setToken(res.token)
    localStorage.setItem(USER_KEY, JSON.stringify(res.user))
    setTokenState(res.token)
    setUser(res.user)
  }, [])

  const register = useCallback(async (data: RegisterData) => {
    const res = await apiRegister(data)
    setToken(res.token)
    localStorage.setItem(USER_KEY, JSON.stringify(res.user))
    setTokenState(res.token)
    setUser(res.user)
  }, [])

  const logout = useCallback(() => {
    clearToken()
    localStorage.removeItem(USER_KEY)
    setTokenState(null)
    setUser(null)
    window.location.href = "/login"
  }, [])

  return {
    user,
    token,
    login,
    register,
    logout,
    isAuthenticated: !!token,
  }
}
