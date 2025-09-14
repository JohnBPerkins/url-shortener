'use client';

import { UrlShortener } from '@/components/UrlShortener';
import { MetricsDashboard } from '@/components/MetricsDashboard';

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-blue-900 to-indigo-900">
      <div className="absolute inset-0 bg-[url('data:image/svg+xml,%3Csvg width="60" height="60" viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg"%3E%3Cg fill="none" fill-rule="evenodd"%3E%3Cg fill="%23ffffff" fill-opacity="0.05"%3E%3Ccircle cx="30" cy="30" r="2"/%3E%3C/g%3E%3C/g%3E%3C/svg%3E')] opacity-40"></div>

      <div className="relative z-10">
        <header className="bg-white/5 backdrop-blur-xl border-b border-white/10 px-4 sm:px-8 py-6">
          <div className="max-w-6xl mx-auto flex flex-col sm:flex-row justify-between items-center gap-4">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-500 rounded-xl flex items-center justify-center">
                <span className="text-white text-xl font-bold">🔗</span>
              </div>
              <h1 className="text-2xl font-bold text-white">
                URL Shortener
              </h1>
            </div>
            <nav className="flex gap-4">
              <a
                href="https://github.com/JohnBPerkins/url-shortener"
                target="_blank"
                rel="noopener noreferrer"
                className="text-white/80 hover:text-white no-underline px-5 py-2.5 rounded-xl transition-all duration-300 bg-white/10 hover:bg-white/20 border border-white/20 flex items-center gap-2 hover:-translate-y-0.5 hover:shadow-xl backdrop-blur-sm"
              >
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.374 0 0 5.373 0 12 0 17.302 3.438 21.8 8.207 23.387c.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.300 24 12c0-6.627-5.373-12-12-12z"/>
                </svg>
                GitHub
              </a>
            </nav>
          </div>
        </header>

        <main className="max-w-6xl mx-auto px-4 sm:px-8 py-16">
          <div className="text-center mb-16">
            <h2 className="text-5xl sm:text-6xl font-bold text-white mb-6">
              Shorten Your
              <span className="bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent block">
                Long URLs
              </span>
            </h2>
            <p className="text-xl text-white/70 max-w-2xl mx-auto">
              Transform lengthy URLs into clean, shareable links in seconds.
              Perfect for social media, emails, and anywhere space matters.
            </p>
          </div>

          <div className="flex flex-col items-center gap-12">
            <UrlShortener />
            <MetricsDashboard />
          </div>
        </main>
      </div>
    </div>
  );
}