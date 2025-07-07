import { SearchBar } from "../components/SearchBar";

export function Welcome() {
  return (
    <main className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-4xl px-4">
        <SearchBar />
      </div>
    </main>
  );
}
