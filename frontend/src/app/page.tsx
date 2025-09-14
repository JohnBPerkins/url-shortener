'use client';

import { UrlShortener } from '@/components/UrlShortener';
import { MetricsDashboard } from '@/components/MetricsDashboard';

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-500 via-purple-600 to-blue-800">
      <header className="bg-white/10 backdrop-blur-lg border-b border-white/20 px-4 sm:px-8 py-4 flex flex-col sm:flex-row justify-between items-center gap-4 shadow-lg">
        <a href="#" className="text-xl font-bold text-white no-underline">
          🔗 URL Shortener
        </a>
        <nav className="flex gap-6">
          <a
            href="https://github.com/JohnBPerkins/url-shortener"
            target="_blank"
            rel="noopener noreferrer"
            className="text-white no-underline px-4 py-2 rounded-lg transition-all duration-300 bg-white/10 border border-white/20 flex items-center gap-2 hover:bg-white/20 hover:-translate-y-0.5 hover:shadow-lg"
          >
            <span>📱</span>
            GitHub
          </a>
        </nav>
      </header>

      <main className="flex flex-col items-center p-8 gap-8">
        <UrlShortener />
        <MetricsDashboard />
      </main>
    </div>
  );
}