// Real curriculum data will be loaded from blockchain via Curriculum module
// This file is for utility functions only - no more mock data

// Prerequisites will be loaded from blockchain and verified via CosmWasm contracts

// Utility function to enhance curriculum with student progress from blockchain NFTs
export const addGradesToCurriculum = (curriculum: any[], studentNFTs: any[]) => {
  return curriculum.map(subject => {
    const nft = studentNFTs.find(nft => nft.subjectId === subject.id)
    return {
      ...subject,
      grade: nft?.grade,
      completed: !!nft,
      completionDate: nft?.completionDate
    }
  })
}
