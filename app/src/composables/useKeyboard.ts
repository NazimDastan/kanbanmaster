import { onMounted, onUnmounted } from 'vue'

type KeyHandler = () => void

interface KeyBinding {
  key: string
  ctrl?: boolean
  handler: KeyHandler
}

export function useKeyboard(bindings: KeyBinding[]) {
  function handleKeyDown(e: KeyboardEvent) {
    // Don't trigger when typing in input/textarea
    const target = e.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) return

    for (const binding of bindings) {
      const ctrlMatch = binding.ctrl ? (e.ctrlKey || e.metaKey) : true
      if (e.key.toLowerCase() === binding.key.toLowerCase() && ctrlMatch) {
        e.preventDefault()
        binding.handler()
        return
      }
    }
  }

  onMounted(() => window.addEventListener('keydown', handleKeyDown))
  onUnmounted(() => window.removeEventListener('keydown', handleKeyDown))
}
