// components/ui/NFTCard.tsx - Reusable component for NFTs
'use client';
import { useState } from 'react';
import { AcademicNFT, Subject } from '../../services/blockchain';

interface NFTCardProps {
  nft: AcademicNFT;
  subject?: Subject;
  onClick?: () => void;
  isSelected?: boolean;
  showActions?: boolean;
  onViewOnBlockchain?: (nftHash: string) => void;
}

export function NFTCard({ 
  nft, 
  subject, 
  onClick, 
  isSelected = false, 
  showActions = true,
  onViewOnBlockchain 
}: NFTCardProps) {
  const [isHovered, setIsHovered] = useState(false);

  const handleViewOnBlockchain = () => {
    if (onViewOnBlockchain) {
      onViewOnBlockchain(nft.nftHash);
    } else {
      // Default explorer URL
      const explorerUrl = process.env.NEXT_PUBLIC_EXPLORER_URL || 'https://explorer.cosmos.network';
      window.open(`${explorerUrl}/tx/${nft.nftHash}`, '_blank');
    }
  };

  return (
    <div 
      className={`
        bg-white rounded-lg shadow-lg p-4 transition-all duration-200 
        ${onClick ? 'cursor-pointer hover:scale-105' : ''}
        ${isSelected ? 'ring-2 ring-blue-500' : ''}
        ${isHovered ? 'shadow-xl' : ''}
      `}
      onClick={onClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* NFT Header */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center space-x-2">
          <span className="text-xs font-mono bg-gray-100 px-2 py-1 rounded">
            {subject?.code || nft.metadata.subject.split(' ')[0]}
          </span>
          <span className="text-xs text-gray-500">
            #{nft.id.slice(-6)}
          </span>
        </div>
        <div className="flex items-center space-x-1">
          <span className="text-lg">✅</span>
          {isHovered && (
            <span className="text-xs text-green-600 font-semibold">NFT</span>
          )}
        </div>
      </div>

      {/* Subject title */}
      <h3 className="font-semibold text-gray-800 mb-2 line-clamp-2">
        {subject?.name || nft.metadata.subject}
      </h3>

      {/* Academic information */}
      <div className="space-y-1 text-sm text-gray-600 mb-3">
        <div className="flex items-center justify-between">
          <span>📚 Credits:</span>
          <span className="font-semibold">{nft.metadata.credits}</span>
        </div>
        <div className="flex items-center justify-between">
          <span>🏛️ Institution:</span>
          <span className="font-semibold text-xs">{nft.metadata.institution}</span>
        </div>
        <div className="flex items-center justify-between">
          <span>📅 Completion:</span>
          <span className="font-semibold text-xs">
            {new Date(nft.completionDate).toLocaleDateString('en-US')}
          </span>
        </div>
      </div>

      {/* NFT Hash (truncated) */}
      <div className="bg-gray-50 rounded p-2 mb-3">
        <div className="text-xs text-gray-500 mb-1">NFT Hash:</div>
        <div className="font-mono text-xs text-gray-700 break-all">
          {nft.nftHash.slice(0, 20)}...{nft.nftHash.slice(-10)}
        </div>
      </div>

      {/* Status badge */}
      <div className="flex items-center justify-between mb-3">
        <span className="px-2 py-1 bg-green-100 text-green-700 rounded-full text-xs font-semibold">
          🏆 Certificate Issued
        </span>
        <span className="text-xs text-gray-500">
          ID: {nft.id.slice(-8)}
        </span>
      </div>

      {/* Actions */}
      {showActions && (
        <div className="space-y-2">
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleViewOnBlockchain();
            }}
            className="w-full bg-blue-500 text-white px-3 py-2 rounded-lg hover:bg-blue-600 transition-colors text-sm font-semibold flex items-center justify-center space-x-2"
          >
            <span>🔗</span>
            <span>View on Blockchain</span>
          </button>
          
          <button
            onClick={(e) => {
              e.stopPropagation();
              // Implement certificate download
              alert('Download feature in development');
            }}
            className="w-full bg-gray-100 text-gray-700 px-3 py-2 rounded-lg hover:bg-gray-200 transition-colors text-sm font-semibold flex items-center justify-center space-x-2"
          >
            <span>📄</span>
            <span>Download Certificate</span>
          </button>
        </div>
      )}

      {/* Hover animation */}
      {isHovered && (
        <div className="absolute inset-0 bg-blue-500/5 rounded-lg pointer-events-none" />
      )}
    </div>
  );
}

// Component for NFT list
interface NFTGridProps {
  nfts: AcademicNFT[];
  subjects?: Subject[];
  onSelectNFT?: (nft: AcademicNFT) => void;
  selectedNFTId?: string;
  loading?: boolean;
  emptyMessage?: string;
}

export function NFTGrid({ 
  nfts, 
  subjects = [], 
  onSelectNFT, 
  selectedNFTId,
  loading = false,
  emptyMessage = "No NFTs found"
}: NFTGridProps) {
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {Array.from({ length: 6 }).map((_, index) => (
          <div key={index} className="bg-white rounded-lg shadow-lg p-4 animate-pulse">
            <div className="flex items-center justify-between mb-3">
              <div className="h-6 bg-gray-200 rounded w-16"></div>
              <div className="h-6 bg-gray-200 rounded w-6"></div>
            </div>
            <div className="h-4 bg-gray-200 rounded mb-2"></div>
            <div className="space-y-2">
              <div className="h-3 bg-gray-200 rounded"></div>
              <div className="h-3 bg-gray-200 rounded"></div>
              <div className="h-3 bg-gray-200 rounded"></div>
            </div>
            <div className="h-8 bg-gray-200 rounded mt-4"></div>
          </div>
        ))}
      </div>
    );
  }

  if (nfts.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-6xl mb-4">🏆</div>
        <h3 className="text-xl font-semibold text-gray-800 mb-2">No Academic NFTs</h3>
        <p className="text-gray-600">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {nfts.map((nft, index) => {
        const subject = subjects.find(s => s.id === nft.subjectId);
        return (
          <NFTCard
            key={`nft-${nft.id || index}`} // Use fallback to avoid duplicate keys
            nft={nft}
            subject={subject}
            onClick={() => onSelectNFT?.(nft)}
            isSelected={selectedNFTId === nft.id}
          />
        );
      })}
    </div>
  );
}

// Component for NFT statistics
interface NFTStatsProps {
  nfts: AcademicNFT[];
  totalCreditsRequired?: number;
}

export function NFTStats({ nfts, totalCreditsRequired = 240 }: NFTStatsProps) {
  const totalCredits = nfts.reduce((sum, nft) => sum + nft.metadata.credits, 0);
  const progress = (totalCredits / totalCreditsRequired) * 100;
  
  const statsByInstitution = nfts.reduce((acc, nft) => {
    acc[nft.metadata.institution] = (acc[nft.metadata.institution] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  return (
    <div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl shadow-lg p-6 text-white">
      <h2 className="text-xl font-semibold mb-4">📊 NFT Statistics</h2>
      
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
        <div>
          <div className="text-3xl font-bold">{nfts.length}</div>
          <div className="text-blue-100">Total NFTs</div>
        </div>
        <div>
          <div className="text-3xl font-bold">{totalCredits}</div>
          <div className="text-blue-100">Credits Earned</div>
        </div>
        <div>
          <div className="text-3xl font-bold">{Math.round(progress)}%</div>
          <div className="text-blue-100">Progress</div>
        </div>
        <div>
          <div className="text-3xl font-bold">{Object.keys(statsByInstitution).length}</div>
          <div className="text-blue-100">Institutions</div>
        </div>
      </div>

      {/* Progress bar */}
      <div className="mb-4">
        <div className="flex justify-between text-sm mb-2">
          <span>Course progress</span>
          <span>{totalCredits}/{totalCreditsRequired} credits</span>
        </div>
        <div className="w-full bg-white/20 rounded-full h-3">
          <div 
            className="bg-white h-3 rounded-full transition-all duration-1000"
            style={{width: `${Math.min(progress, 100)}%`}}
          ></div>
        </div>
      </div>

      {/* Distribution by institution */}
      {Object.keys(statsByInstitution).length > 1 && (
        <div>
          <div className="text-sm text-blue-100 mb-2">NFTs by institution:</div>
          <div className="flex flex-wrap gap-2">
            {Object.entries(statsByInstitution).map(([institution, count]) => (
              <span key={institution} className="bg-white/20 px-2 py-1 rounded text-xs">
                {institution}: {count}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}