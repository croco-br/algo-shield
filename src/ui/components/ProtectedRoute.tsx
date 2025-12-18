'use client';

import { useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { useAuth } from '@/lib/stores/auth';

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
	const { user, loading } = useAuth();
	const router = useRouter();
	const pathname = usePathname();

	useEffect(() => {
		const publicRoutes = ['/login'];
		const isPublicRoute = publicRoutes.some((route) => pathname.startsWith(route));

		if (!loading && !user && !isPublicRoute) {
			router.push('/login');
		}
		if (!loading && user && pathname.startsWith('/login')) {
			router.push('/');
		}
	}, [user, loading, pathname, router]);

	if (loading && !user && !pathname.startsWith('/login')) {
		return (
			<div className="flex items-center justify-center min-h-[50vh]">
				<div className="w-10 h-10 border-4 border-gray-200 border-t-indigo-600 rounded-full animate-spin"></div>
			</div>
		);
	}

	return <>{children}</>;
}
