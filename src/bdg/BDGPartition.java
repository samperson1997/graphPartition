/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Queue;
import java.util.Set;
import java.util.TreeSet;

public class BDGPartition {

  /**
   * index: block index (color), which represents different blocks
   */
  private TreeSet<Block> blocks;

  /**
   * index: block index, which represents neighbors of the block
   */
  private List<List<Integer>> blockNeighborList;

  /**
   * index: worker index, which represents the partition result
   */
  private List<List<Block>> workerBlocks;

  private void BDG(List<Vertex> vertices, int workerNum, int blockNum, int step) {
    // init blocks
    init(blockNum, workerNum);

    // cut an input graph into fine-grained blocks
    bfs(vertices, blockNum, step);

    // apply deterministic greedy algorithm
    workerBlocks.get(0).add(blocks.first());
    deterministicGreedy(vertices.size(), blockNum);
  }

  private void init(int blockNum, int workerNum) {
    blocks = new TreeSet<>();
    blockNeighborList = new ArrayList<>();
    for (int i = 0; i < blockNum; i++) {
      blockNeighborList.add(new ArrayList<>());
    }
    workerBlocks = new ArrayList<>();
    for (int i = 0; i < workerNum; i++) {
      workerBlocks.add(new ArrayList<>());
    }
  }

  /**
   * @param vertexes graph
   * @param blockNum number of blocks (config)
   * @param step number of steps taken by BFS (config)
   */
  private void bfs(List<Vertex> vertexes, int blockNum, int step) {
    Queue<Vertex> queue = new LinkedList<>();
    Set<Vertex> srcSet = new HashSet<>();
    int idx = 0;
    while (srcSet.size() < blockNum) {
      Vertex src = vertexes.get(idx);
      if (!srcSet.contains(src)) {
        src.setColor(idx);
        queue.add(src);
        srcSet.add(src);

        // put src node into block map
        Set<Vertex> blockSet = new HashSet<>();
        blockSet.add(src);
        blocks.add(new Block(idx, blockSet));
      }
      idx++;
    }

    while (!queue.isEmpty()) {
      int size = queue.size();
      for (int i = 0; i < size; i++) {
        Vertex p = queue.poll();
        for (Vertex neighbor : p.getNeighbors()) {
          if (neighbor.getColor() == -1) {
            neighbor.setColor(p.getColor());
            for (Block block : blocks) {
              if (block.getId() == p.getColor()) {
                block.getVertices().add(neighbor);
                break;
              }
            }
            queue.offer(neighbor);
          } else {
            blockNeighborList.get(p.getColor()).add(neighbor.getColor());
          }
        }
      }
    }

    // running CC finding algorithm (like Hash-Min) on uncolored vertices
  }

  /**
   * deterministic greedy
   *
   * @param capacity size of vertices
   * @param blockNum number of blocks
   */
  private void deterministicGreedy(int capacity, int blockNum) {
    for (int i = 1; i < blockNum; i++) {
      Block block = blocks.pollFirst();
      if (block == null) {
        return;
      }

      // for each 1-hop neighbor blocks of this block
      List<Integer> neighborBlocks = blockNeighborList.get(block.getId());
      Set<Vertex> bSet = new HashSet<>();
      for (int neighborBlock : neighborBlocks) {
        bSet.addAll(getVerticesByBlockId(neighborBlock));
      }

      int j = 0;
      Set<Vertex> pSet = new HashSet<>();

      // for each worker i, calculate j
      for (List<Block> blocksInWorker : workerBlocks) {
        for (Block blockInWorker : blocksInWorker) {
          pSet.addAll(blockInWorker.getVertices());
        }
        bSet.retainAll(pSet);
        j = Math.max(j, bSet.size() * (1 - pSet.size() / capacity));
      }

      workerBlocks.get(j).add(block);
    }
  }

  private Set<Vertex> getVerticesByBlockId(int blockId) {
    for (Block block : blocks) {
      if (block.getId() == blockId) {
        return block.getVertices();
      }
    }
    return Collections.emptySet();
  }

  private class Block implements Comparable {

    private int id;
    private Set<Vertex> vertices;

    Block(int id, Set<Vertex> vertices) {
      this.id = id;
      this.vertices = vertices;
    }

    int getId() {
      return id;
    }

    Set<Vertex> getVertices() {
      return vertices;
    }

    @Override
    public int compareTo(Object o) {
      return Integer.compare(vertices.size(), ((Block) o).vertices.size());
    }
  }
}
