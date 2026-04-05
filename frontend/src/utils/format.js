export const truncateAddress = (addr, start = 6, end = 4) => {
  if (!addr) return ''
  return addr.slice(0, start) + '...' + addr.slice(-end)
}

export const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleString()
}

export const formatSOL = (value) => {
  return (value || 0).toFixed(4)
}