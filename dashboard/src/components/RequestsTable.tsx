import type { Request } from '../types'

interface Props {
  requests: Request[]
}

function statusColor(code: number): string {
  if (code >= 500) return 'text-red-400'
  if (code >= 400) return 'text-yellow-400'
  return 'text-green-400'
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString()
}

export default function RequestsTable({ requests }: Props) {
  return (
    <div className="mb-8">
      <h2 className="text-white text-lg font-semibold mb-3">Requests</h2>
      <div className="border border-gray-800 rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-900">
            <tr>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">ID</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Model</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Status</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Latency</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Time</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {requests.length === 0 && (
              <tr>
                <td colSpan={5} className="text-gray-500 text-center px-4 py-8">
                  No requests yet
                </td>
              </tr>
            )}
            {requests.map((r) => (
              <tr key={r.id} className="bg-gray-950 hover:bg-gray-900 transition-colors">
                <td className="px-4 py-3 text-gray-400 font-mono text-xs">
                  {r.id.slice(0, 8)}...
                </td>
                <td className="px-4 py-3 text-white">{r.model || '—'}</td>
                <td className={`px-4 py-3 font-mono ${statusColor(r.status_code)}`}>
                  {r.status_code}
                </td>
                <td className="px-4 py-3 text-gray-300">{r.latency_ms}ms</td>
                <td className="px-4 py-3 text-gray-400">{formatDate(r.created_at)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}