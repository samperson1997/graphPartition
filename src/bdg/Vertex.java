import java.util.List;

public class Vertex {

  private String name;
  private List<Vertex> neighbors;
  private int color;

  public Vertex(String name, List<Vertex> neighbors) {
    this.name = name;
    this.neighbors = neighbors;
    this.color = -1;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public List<Vertex> getNeighbors() {
    return neighbors;
  }

  public void setNeighbors(List<Vertex> neighbors) {
    this.neighbors = neighbors;
  }

  public int getColor() {
    return color;
  }

  public void setColor(int color) {
    this.color = color;
  }
}
