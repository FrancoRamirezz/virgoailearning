"use client";

import React, { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Eye, EyeOff, AlertCircle, BookOpen } from "lucide-react";
import { FaGoogle, FaFacebookF, FaApple } from "react-icons/fa";
import { signup } from "@/lib/auth";

interface SignUpFormData {
    name: string;
    email: string;
    password: string;
    confirmPassword: string;
    rememberMe: boolean;
}

export default function SignUpPage() {
  const router = useRouter();
  const [formData, setFormData] = useState<SignUpFormData>({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
    rememberMe: false,
  });

  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL;

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ): void => {
    const { name, value, type, checked } = e.target;

    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }));

    if (error) setError(null);
  };

  const handleSubmit = async ( e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      if (formData.password !== formData.confirmPassword) {
        throw new Error("Passwords do not match.");
      }

      await signup({
        role: "student",
        name: formData.name,
        email: formData.email,
        password: formData.password,
        rememberMe: formData.rememberMe,
      });

      router.replace("/dashboard");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Something went wrong. Please try again."
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="w-full lg:w-1/2 flex items-center justify-center p-8">
      <div className="w-full max-w-md">
        {/* Logo */}
        <div className="mb-12">
          <Link href="/" className="flex items-center gap-3 mb-8">
            <Image src="/logo.png" alt="VirgoLearning Logo" width={70} height={70} />
            {/* <div className="w-10 h-10 bg-gradient-to-r from-red-600 to-blue-600 rounded-lg flex items-center justify-center">
              <BookOpen className="w-6 h-6 text-white" />
            </div> */}
            <span className="text-2xl font-bold text-gray-900">CitizenshipPrep</span>
          </Link>
          <h1 className="text-3xl font-semibold text-gray-900 mb-2">
            Create your account
          </h1>
          <p className="text-gray-600">
            Start your learning journey
          </p>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-start gap-3">
              <AlertCircle className="w-5 h-5 text-red-600 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-red-700">{error}</p>
            </div>
          </div>
        )}

        {/* Login Form */}
        <form onSubmit={handleSubmit} className="space-y-5">
            {/* Name Field */}
            <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
                Name
                </label>
                <input
                    id="name"
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleInputChange}
                    required
                    autoComplete="name"
                    className="w-full px-3 py-2.5 text-base border border-gray-300 rounded-md focus:border-red-600 focus:ring-1 focus:ring-red-600 transition-colors outline-none"
                    placeholder="Full name"
                />
            </div>

          {/* Email Field */}
            <div>
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                Email
                </label>
                <input
                    id="email"
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleInputChange}
                    required
                    autoComplete="email"
                    className="w-full px-3 py-2.5 text-base border border-gray-300 rounded-md focus:border-red-600 focus:ring-1 focus:ring-red-600 transition-colors outline-none"
                    placeholder="Email address"
                />
            </div>

          {/* Password Field */}
            <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
                Password
                </label>
                <div className="relative">
                <input
                    id="password"
                    type={showPassword ? 'text' : 'password'}
                    name="password"
                    value={formData.password}
                    onChange={handleInputChange}
                    required
                    autoComplete="current-password"
                    className="w-full px-3 py-2.5 text-base border border-gray-300 rounded-md focus:border-red-600 focus:ring-1 focus:ring-red-600 transition-colors outline-none pr-10"
                    placeholder="Password"
                />
                <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                    aria-label={showPassword ? "Hide password" : "Show password"}
                >
                    {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                    ) : (
                    <Eye className="w-5 h-5" />
                    )}
                </button>
                </div>
            </div>

            {/* Confirm Password Field */}
            <div>
                <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-2">
                Confirm Password
                </label>
                <div className="relative">
                <input
                    id="confirmPassword"
                    type={showPassword ? 'text' : 'password'}
                    name="confirmPassword"
                    value={formData.confirmPassword}
                    onChange={handleInputChange}
                    required
                    autoComplete="current-password"
                    className="w-full px-3 py-2.5 text-base border border-gray-300 rounded-md focus:border-red-600 focus:ring-1 focus:ring-red-600 transition-colors outline-none pr-10"
                    placeholder="Password"
                />
                <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showPassword)}
                    className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                    aria-label={showPassword ? "Hide password" : "Show password"}
                >
                    {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                    ) : (
                    <Eye className="w-5 h-5" />
                    )}
                </button>
                </div>
            </div>

          {/* Remember Me*/}
          <div className="flex items-center justify-between text-sm">
            <label className="flex items-center cursor-pointer group">
              <input
                type="checkbox"
                name="rememberMe"
                checked={formData.rememberMe}
                onChange={handleInputChange}
                className="w-4 h-4 text-red-600 border-gray-300 rounded focus:ring-2 focus:ring-red-500 cursor-pointer"
              />
              <span className="ml-2 text-gray-700 group-hover:text-gray-900">
                Keep me signed in
              </span>
            </label>
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            disabled={isSubmitting}
            className={`w-full py-2.5 px-4 rounded-md font-medium text-white transition-all duration-150 ${
              isSubmitting
                ? 'bg-gray-400 cursor-not-allowed'
                : 'btn-secondary'
                // : 'bg-red-600 hover:bg-red-700 active:bg-red-800'
            }`}
          >
            {isSubmitting ? (
              <span className="flex items-center justify-center gap-2">
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                Creating Account...
              </span>
            ) : (
              'Create Account'
            )}
          </button>
        </form>

        {/* Divider */}
        <div className="relative my-8">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-200"></div>
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-slate-50 text-gray-500">or</span>
          </div>
        </div>

        {/* OAuth Login Options */}
        <div className="flex justify-center gap-4 mb-8">
          {/* Google */}
          <button
            type="button"
            onClick={() => (window.location.href = `${API_BASE}/auth/google/login`)}
            aria-label="Sign in with Google"
            className="flex items-center justify-center w-12 h-12 rounded-full bg-white border border-gray-200 shadow-sm hover:shadow-md transition hover:scale-105"
          >
            <FaGoogle className="text-red-500 text-xl" />
          </button>

          {/* Facebook */}
          <button
            type="button"
            onClick={() => (window.location.href = `${API_BASE}/auth/facebook/login`)}
            aria-label="Sign in with Facebook"
            className="flex items-center justify-center w-12 h-12 rounded-full bg-white border border-gray-200 shadow-sm hover:shadow-md transition hover:scale-105"
          >
            <FaFacebookF className="text-blue-600 text-xl" />
          </button>
        </div>

        {/* Sign Up Link */}
        <div className="text-center">
          <p className="text-sm text-gray-600">
            Already have an account?{' '}
            <Link href="/login" className="text-red-600 hover:text-red-700 font-medium">
              Login
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

