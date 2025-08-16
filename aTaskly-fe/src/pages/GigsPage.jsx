import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useSearchParams } from 'react-router-dom';
import GigCard from '../components/GigCard';
import { apiGet } from '../utils/api';
import './GigsPage.css';

// Debounce utility function
const debounce = (func, delay) => {
  let timeout;
  return function(...args) {
    const context = this;
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(context, args), delay);
  };
};

const GigsPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [gigs, setGigs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [nextSearchAfter, setNextSearchAfter] = useState(null);
  const [loadingMore, setLoadingMore] = useState(false);

  const nextSearchAfterRef = useRef(nextSearchAfter);
  nextSearchAfterRef.current = nextSearchAfter;

  const [searchTerm, setSearchTerm] = useState(searchParams.get('search_term') || '');

  const [filters, setFilters] = useState({
    category: searchParams.get('category') || '',
    minPrice: searchParams.get('minPrice') || '',
    maxPrice: searchParams.get('maxPrice') || '',
  });

  const observer = useRef();
  const fetchGigs = useCallback(async (isLoadMore = false) => {
    if (isLoadMore) {
      setLoadingMore(true);
    } else {
      setLoading(true);
      setGigs([]); // Reset gigs on new search
    }
    setError('');

    try {
      const queryParams = new URLSearchParams();

      // Search Term
      if (searchTerm) {
        queryParams.set('search_term', searchTerm);
      }

      // Min Price
      if (filters.minPrice) {
        queryParams.set('min_price', filters.minPrice);
      }

      // Max Price
      if (filters.maxPrice) {
        queryParams.set('max_price', filters.maxPrice);
      }

      // Category IDs (assuming single for now, can be extended for multiple)
      if (filters.category) {
        queryParams.set('category_ids', filters.category);
      }

      // Last Gig ID for infinite scroll
      if (isLoadMore && nextSearchAfterRef.current) {
        queryParams.set('last_gig_id', nextSearchAfterRef.current);
      }

      const endpoint = `/gigs/search?${queryParams.toString()}`;
      const data = await apiGet(endpoint);
      
      const newGigs = data; // data is directly []model.SearchGigDTO from backend
      setGigs(prevGigs => isLoadMore ? [...prevGigs, ...newGigs] : newGigs);

      // Determine nextSearchAfter from the last gig in the new batch
      if (newGigs && newGigs.length > 0 && newGigs.length === 10) { // Assuming backend returns max 10 gigs per request
          setNextSearchAfter(newGigs[newGigs.length - 1].id);
      } else {
          setNextSearchAfter(null); // No more results or less than limit
      }

    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Failed to fetch gigs.';
      setError(errorMessage);
    } finally {
      setLoading(false);
      setLoadingMore(false);
    }
  }, [searchTerm, filters]);

  useEffect(() => {
    setSearchTerm(searchParams.get('search_term') || '');
    setFilters({
      category: searchParams.get('category') || '',
      minPrice: searchParams.get('minPrice') || '',
      maxPrice: searchParams.get('maxPrice') || '',
    });
  }, [searchParams]);

  useEffect(() => {
    setNextSearchAfter(null); // Reset for new search
    fetchGigs();
  }, [searchTerm, filters, fetchGigs]);

  const lastGigElementRef = useCallback(node => {
    if (loadingMore) return;
    if (observer.current) observer.current.disconnect();
    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && nextSearchAfterRef.current) {
        fetchGigs(true);
      }
    });
    if (node) observer.current.observe(node);
  }, [loadingMore, fetchGigs]);

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters(prev => ({ ...prev, [name]: value }));
  };

  const applyFilters = () => {
    const newSearchParams = new URLSearchParams(searchParams);
    
    // Handle search term
    if (searchTerm) {
      newSearchParams.set('search_term', searchTerm);
    } else {
      newSearchParams.delete('search_term');
    }

    // Handle other filters
    Object.keys(filters).forEach(key => {
      if (filters[key]) {
        newSearchParams.set(key, filters[key]);
      } else {
        newSearchParams.delete(key);
      }
    });
    setSearchParams(newSearchParams);
  };

  if (loading) {
    return <div className="loading-state">Finding services for you...</div>;
  }

  if (error) {
    return <div className="error-message">{error}</div>;
  }

  return (
    <div className="gigs-page">
      <div className="container">
        <div className="page-header">
          <h1>Results for "{searchParams.get('search_term')}"</h1>
          <div className="filters">
            <input
              type="text"
              name="searchTerm"
              placeholder="Search for services..."
              value={searchTerm}
              onChange={debounce((e) => setSearchTerm(e.target.value), 500)}
            />
            <input
              type="text"
              name="category"
              placeholder="Category"
              value={filters.category}
              onChange={handleFilterChange}
            />
            <input
              type="number"
              name="minPrice"
              placeholder="Min Price"
              value={filters.minPrice}
              onChange={handleFilterChange}
            />
            <input
              type="number"
              name="maxPrice"
              placeholder="Max Price"
              value={filters.maxPrice}
              onChange={handleFilterChange}
            />
            <button onClick={applyFilters}>Apply Filters</button>
          </div>
        </div>
        <div className="gig-list">
          {gigs.length > 0 ? (
            gigs.map((item, index) => <GigCard key={`${item.id}-${index}`} gig={item} />)
          ) : (
            <p>No services found for your search.</p>
          )}
        </div>
        {nextSearchAfter && (
          <div ref={lastGigElementRef} className="load-more-trigger">
            {loadingMore && <div className="loading-state">Loading more...</div>}
          </div>
        )}
      </div>
    </div>
  );
};

export default GigsPage; 