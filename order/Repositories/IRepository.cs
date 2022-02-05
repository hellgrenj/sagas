namespace model.Repositories;

public interface IRepository<T>
{
    Task<T> GetAsync(int id);
    Task<int> SaveAsync(T entity);
}