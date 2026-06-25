import type { Anomaly } from '../types'

interface Props {
  anomalies: Anomaly[]
}

function ruleColor(rule: string): string {
  if (rule === 'upstream_error') return 'text-red-400 bg-red-950'
  if (rule === 'client_error') return 'text-yellow-400 bg-yellow-950'
  if (rule === 'high_latency') return 'text-orange-400 bg-orange-950'
  if (rule === 'high_token_usage') return 'text-purple-400 bg-purple-950'
  return 'text-gray-400 bg-gray-800'
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString()
}

export default function AnomaliesTable({ anomalies }: Props) {
  return (
    <div>
      <h2 className="text-white text-lg font-semibold mb-3">Anomalies</h2>
      <div className="border border-gray-800 rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-900">
            <tr>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Rule</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Detail</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Request ID</th>
              <th className="text-left text-gray-400 px-4 py-3 font-medium">Time</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-800">
            {anomalies.length === 0 && (
              <tr>
                <td colSpan={4} className="text-gray-500 text-center px-4 py-8">
                  No anomalies detected
                </td>
              </tr>
            )}
            {anomalies.map((a) => (
              <tr key={a.id} className="bg-gray-950 hover:bg-gray-900 transition-colors">
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs font-mono ${ruleColor(a.rule)}`}>
                    {a.rule}
                  </span>
                </td>
                <td className="px-4 py-3 text-gray-300">{a.detail}</td>
                <td className="px-4 py-3 text-gray-400 font-mono text-xs">
                  {a.request_id.slice(0, 8)}...
                </td>
                <td className="px-4 py-3 text-gray-400">{formatDate(a.created_at)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}