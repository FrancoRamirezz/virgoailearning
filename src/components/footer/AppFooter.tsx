export default function AppFooter() {
  return (
    <footer className="border-t border-gray-200 bg-(--color-blue-600)">
      <div className="px-4 py-4 mx-auto max-w-7xl flex flex-col sm:flex-row items-center justify-between text-sm text-gray-500">
        
        <p className="text-gray-400">&copy; 2025 VirgoLearning. All rights reserved.</p>

        <div className="flex space-x-4 mt-2 sm:mt-0">
          <a href="/help" className="hover:text-(--color-gold)">Help</a>
          <a href="/faq" className="hover:text-(--color-gold)">FAQ</a>
          <a href="/privacy" className="hover:text-(--color-gold)">Privacy</a>
          <a href="/terms" className="hover:text-(--color-gold)">Terms</a>
        </div>
      </div>
    </footer>
  );
}
