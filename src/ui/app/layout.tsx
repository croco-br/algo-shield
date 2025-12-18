import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/lib/stores/auth";
import Header from "@/components/Header";
import ProtectedRoute from "@/components/ProtectedRoute";

export const metadata: Metadata = {
  title: "AlgoShield",
  description: "Fraud Detection & Anti-Money Laundering",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <ProtectedRoute>
            <div className="min-h-screen bg-gray-50">
              <Header />
              <main className="py-8">{children}</main>
            </div>
          </ProtectedRoute>
        </AuthProvider>
      </body>
    </html>
  );
}
