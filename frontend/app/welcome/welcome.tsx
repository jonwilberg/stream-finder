import { SearchBar } from "../components/SearchBar";

export function Welcome() {
  return (
    <main className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-4xl px-4">
        <div className="mb-12 text-center">
          <img 
            src="/logo.png" 
            alt="Stream Finder Logo" 
            className="mx-auto h-48 w-auto"
          />
        </div>
        <SearchBar />
      </div>
    </main>
  );
}
