'use client';

import { useState, useEffect, useRef } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/lib/stores/auth';

export default function Header() {
	const { user, logout } = useAuth();
	const pathname = usePathname();
	const router = useRouter();
	const [showUserMenu, setShowUserMenu] = useState(false);
	const menuRef = useRef<HTMLDivElement>(null);

	const isAdmin = user?.roles?.some((role) => role.name === 'admin') || false;

	useEffect(() => {
		function handleClickOutside(event: MouseEvent) {
			if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
				setShowUserMenu(false);
			}
		}

		document.addEventListener('click', handleClickOutside);
		return () => document.removeEventListener('click', handleClickOutside);
	}, []);

	const handleLogout = async () => {
		await logout();
		router.push('/login');
	};

	if (!user || pathname.startsWith('/login')) {
		return null;
	}

	return (
		<header className="bg-white border-b border-gray-200 shadow-sm">
			<div className="max-w-7xl mx-auto px-8">
				<nav className="py-6 flex items-center justify-between gap-8">
					<div className="flex items-center gap-3">
						<img src="/gopher.png" alt="AlgoShield" className="w-10 h-10 object-contain" />
						<div>
							<h1 className="text-2xl font-semibold text-gray-900">AlgoShield</h1>
							<p className="text-sm text-gray-500">Fraud Detection & Anti-Money Laundering</p>
						</div>
					</div>

					<div className="flex gap-2 flex-1 justify-center">
						<Link
							href="/"
							className={`px-6 py-3 rounded-md text-sm font-medium transition-all ${
								pathname === '/'
									? 'text-indigo-600 bg-indigo-50'
									: 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'
							}`}
						>
							Rules
						</Link>
						<Link
							href="/synthetic-test"
							className={`px-6 py-3 rounded-md text-sm font-medium transition-all ${
								pathname === '/synthetic-test'
									? 'text-indigo-600 bg-indigo-50'
									: 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'
							}`}
						>
							Synthetic Test
						</Link>
						{isAdmin && (
							<Link
								href="/permissions"
								className={`px-6 py-3 rounded-md text-sm font-medium transition-all ${
									pathname === '/permissions'
										? 'text-indigo-600 bg-indigo-50'
										: 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'
								}`}
							>
								Permissions
							</Link>
						)}
					</div>

					<div className="relative" ref={menuRef}>
						<button
							onClick={() => setShowUserMenu(!showUserMenu)}
							className="flex items-center gap-3 px-2 py-2 border border-gray-200 rounded-lg bg-white hover:bg-gray-50 transition-all"
						>
							<div className="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-white flex items-center justify-center font-semibold text-lg overflow-hidden">
								{user.picture_url ? (
									<img src={user.picture_url} alt={user.name} className="w-full h-full object-cover" />
								) : (
									user.name.charAt(0).toUpperCase()
								)}
							</div>
							<div className="text-left">
								<div className="text-sm font-medium text-gray-900">{user.name}</div>
								<div className="text-xs text-gray-500">{user.email}</div>
							</div>
							<svg
								className="w-4 h-4 text-gray-500"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<polyline points="6 9 12 15 18 9"></polyline>
							</svg>
						</button>

						{showUserMenu && (
							<div className="absolute top-full right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg min-w-[240px] z-50">
								<div className="p-4 border-b border-gray-200">
									<div className="font-semibold text-gray-900 mb-1">{user.name}</div>
									<div className="text-sm text-gray-500 mb-2">{user.email}</div>
									<div className="flex flex-wrap gap-2 mt-2">
										{(user.roles || []).map((role) => (
											<span
												key={role.id}
												className="inline-block px-3 py-1 bg-indigo-600 text-white rounded text-xs font-medium"
											>
												{role.name}
											</span>
										))}
									</div>
								</div>
								<button
									onClick={handleLogout}
									className="w-full flex items-center gap-3 px-4 py-3 text-sm text-red-600 hover:bg-gray-50 transition-colors border-t border-gray-200"
								>
									<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
										<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
										<polyline points="16 17 21 12 16 7"></polyline>
										<line x1="21" y1="12" x2="9" y2="12"></line>
									</svg>
									Logout
								</button>
							</div>
						)}
					</div>
				</nav>
			</div>
		</header>
	);
}
