"use client"

import { useQueryState } from "nuqs"
import { Input } from "@/components/ui/input"
import { useDebouncedCallback } from "use-debounce"

export function SearchBar() {
  const [search, setSearch] = useQueryState("search", { defaultValue: "" })

  const handleSearch = useDebouncedCallback((term: string) => {
    setSearch(term || null)
  }, 300)

  return (
    <div className="w-full max-w-sm">
      <Input
        placeholder="Search..."
        defaultValue={search}
        onChange={(e) => handleSearch(e.target.value)}
      />
    </div>
  )
}
