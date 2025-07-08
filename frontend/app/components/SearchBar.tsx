import { useState } from "react";

export function SearchBar() {
  const [searchValue, setSearchValue] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Search:", searchValue);
  };

  return (
    <form onSubmit={handleSubmit} className="w-full max-w-2xl mx-auto px-4 sm:px-0">
      <div className="relative">
        <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
          <svg
            className="h-5 w-5 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
        <input
          type="text"
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
          placeholder="Search for a movie or TV show"
          className="w-full pl-12 pr-20 py-3 sm:py-4 text-base sm:text-lg text-gray-900 bg-white border border-gray-200 rounded-2xl shadow-lg focus:outline-none focus:border-blue-200 focus:ring-0 placeholder-gray-400"
        />
        <button
          type="submit"
          className="absolute right-2 top-1/2 transform -translate-y-1/2 aspect-square h-10 sm:h-12 bg-gray-100 hover:bg-gray-200 rounded-2xl flex items-center justify-center transition-colors duration-200 cursor-pointer"
        >
          <svg
            className="h-5 w-5 text-gray-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M9 5l7 7-7 7"
            />
          </svg>
        </button>
      </div>
    </form>
  );
}